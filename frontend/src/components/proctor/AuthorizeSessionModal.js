import React, { useState, useEffect } from "react";
import axios from 'axios';
import { Container, Modal, Spinner, Button, Form } from "react-bootstrap";
import Timer from "../participant/Timer";
import { toast } from 'react-toastify';

const AuthorizeSessionModal = (props) => {
  const NOT_STARTED = 0;
  const ALREADY_STARTED = 1;
  const ALREADY_FINISHED = 2;

  const { show, handleClose, auth, examSession } = props;
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [state, setState] = useState(NOT_STARTED);

  const fields = [
    {
      label: 'Durasi Maksimal Pengerjaan yang Diperbolehkan (dalam satuan menit)',
      name: 'allowed_duration_minutes',
      placeholder: 'Masukkan durasi pengerjaan yang diperbolehkan untuk peserta ini',
      type: 'float',
      minValue: 1,
      defaultValue: 120,
      step: 1,
    },
  ];

  const [formData, setFormData] = useState(
    fields.reduce((acc, field) => ({ ...acc, [field.name]: field.defaultValue }), {})
  );

  useEffect(() => {
    if (examSession === '') {
      return;
    }
    setLoading(true);
    setError(null);
    setData(null);
    setState(NOT_STARTED);

    axios.get(`${process.env.REACT_APP_BACKEND_URL}/api/v1/proctor/participant-sessions/${examSession}/check`, {
      headers: {
        'Authorization': `Bearer ${auth.token}`
      },
    })
    .then(response => {
      setData(response.data.data);
      if (response.data.data.is_submitted) {
        setState(ALREADY_FINISHED);
      }
      else if (response.data.data.participant.started_at) {
        setState(ALREADY_STARTED);
      } else {
        setState(NOT_STARTED);
      }
      setLoading(false);
    })
    .catch(err => {
      if (err.status === 500) {
        setError('Sedang terjadi masalah pada server. Silakan coba beberapa saat lagi.');
      } else {
        setError(error);
      }
      setLoading(false);
    });
  }, [examSession]);

  const handleAuthorize = () => {
    setLoading(true);
    axios.post(`${process.env.REACT_APP_BACKEND_URL}/api/v1/proctor/participant-sessions/${examSession}/authorize`, {
      allowed_duration_minutes: formData.allowed_duration_minutes,
    }, {
      headers: {
        'Authorization': `Bearer ${auth.token}`
      },
    })
    .then(response => {
      toast.success('Sesi berhasil diizinkan!', {
        position: "top-center",
        autoClose: 3000,
        hideProgressBar: false,
        closeOnClick: true,
        pauseOnHover: true,
        draggable: true,
      });
      handleClose();
    })
    .catch(err => {
      toast.error(`Sedang terjadi masalah pada server. Silakan coba beberapa saat lagi.`, {
        position: "top-center",
        autoClose: 5000,
        hideProgressBar: false,
        closeOnClick: true,
        pauseOnHover: true,
        draggable: true,
      });
      handleClose();
    })
  }

  const handleInputChange = (e) => {
    const { name, type } = e.target;
    if (type === 'checkbox') {
      const { checked } = e.target;
      setFormData({ ...formData, [name]: checked });
    } else {
      const { value } = e.target;
      if (name === 'allowed_duration_minutes' && value === '') {
        return;
      }
      if (name === 'allowed_duration_minutes') {
        setFormData({ ...formData, [name]: Number(value) });
      } else {
        setFormData({ ...formData, [name]: value });
      }
    }
  };

  return (
    <Modal show={show} onHide={handleClose} size="lg" centered>
      <Modal.Header closeButton>
        <Modal.Title>Data Sesi Ujian</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <Container>
          { loading
          ? (
            <Container className="text-center">
              <Spinner animation="border" />
              <p>Mohon tunggu...</p>
            </Container>
          )
          : (
            <>
            {
              error ? (
                <>
                  { error }
                </>
              )
              : (
                <>
                  <Form.Group controlId="examName" className="mt-3">
                    <Form.Label><b>Nama Ujian</b></Form.Label>
                    <Form.Control
                      type="text"
                      value={data.exam.name}
                      autoComplete='off'
                      disabled={true}
                    />
                  </Form.Group>
                  
                  <Form.Group controlId="participantName" className="mt-3">
                    <Form.Label><b>Kode Peserta</b></Form.Label>
                    <Form.Control
                      type="text"
                      value={data.participant.name}
                      autoComplete='off'
                      disabled={true}
                    />
                  </Form.Group>

                  <Form.Group controlId="status" className="mt-3">
                    <Form.Label><b>Status</b></Form.Label>
                    <Form.Control
                      type="text"
                      value={
                        state === NOT_STARTED
                        ? "Belum dimulai"
                        : state === ALREADY_STARTED
                        ? "Sudah dimulai"
                        : "Sudah selesai"
                      }
                      autoComplete='off'
                      disabled={true}
                    />
                  </Form.Group>

                  {
                    state === NOT_STARTED
                    ? (
                      <>
                        {fields.map((field) => (
                          <Form.Group className="my-3" controlId={field.name} key={field.name}>
                            <Form.Label><b>{field.label}{field.required ? <span> (harus diisi)</span> : ''}</b></Form.Label>
                            {
                              field.type === 'float'
                              ? (
                                <Form.Control
                                  type='number'
                                  min={field.minValue}
                                  step={field.step}
                                  name={field.name}
                                  value={formData[field.name]}
                                  onChange={handleInputChange}
                                  required={field.required}
                                  placeholder={`Masukkan ${field.label}`}
                                  autoComplete='off'
                                />
                              )
                              : field.type === 'boolean'
                              ? (
                                <Form.Check
                                  type="switch"
                                  name={field.name}
                                  checked={formData[field.name]}
                                  onChange={handleInputChange}
                                />
                              )
                              : (
                                <Form.Control
                                  type={field.type}
                                  name={field.name}
                                  value={formData[field.name]}
                                  onChange={handleInputChange}
                                  required={field.required}
                                  placeholder={`Masukkan ${field.label}`}
                                  autoComplete='off'
                                />
                              )
                            }
                          </Form.Group>
                        ))}
                      </>
                    )
                    : (
                      <Form.Group controlId="timeLeft" className="mt-3">
                        <Form.Label><b>Waktu Tersisa</b></Form.Label>
                        <Timer startTime={data.participant.started_at} durationMinutes={
                          data.participant.is_submitted ? 0 : data.participant.allowed_duration_minutes
                        }
                        onTimesUp={() => {}}></Timer>
                      </Form.Group>
                    )
                  }
                </>
              )
            }
            </>
          )}
        </Container>
      </Modal.Body>
      <Modal.Footer>
        <Button variant="danger" onClick={handleClose}>
          Batal
        </Button>
        <Button variant="primary" onClick={handleAuthorize} disabled={loading || error || (state === ALREADY_FINISHED)}>
          Izinkan
        </Button>
      </Modal.Footer>
    </Modal>
  )
};

export default AuthorizeSessionModal;