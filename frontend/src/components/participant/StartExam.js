import axios from 'axios';
import React, { useEffect, useState } from 'react';
import { Button, Col, Container, Form, Row, Spinner } from 'react-bootstrap';
import { toast } from 'react-toastify';
import { useNavigate, useParams } from 'react-router-dom';
import {QRCodeSVG} from 'qrcode.react';

const StartExam = () => {
  const { examSerial } = useParams();
  const navigate = useNavigate();
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(true);
  const [name, setName] = useState('');
  const [fetchedExam, setFetchedExam] = useState({});
  const [examSessionSerial, setExamSessionSerial] = useState('');
  const [isSessionAuthorized, setIsSessionAuthorized] = useState(false);

  useEffect(() => {
    axios.get(`${process.env.REACT_APP_BACKEND_URL}/api/v1/exams/${examSerial}`)
    .then(response => {
      if (response.data.data) {
        setFetchedExam(response.data.data);
      } else {
        navigate('/404');
      }
      setLoading(false);
    })
    .catch(error => {
      if (error.status === 404) {
        navigate('/404');
      }
      setError(error.message);
      setLoading(false);
    });
  }, []);

  if (loading) {
    return (
      <Container className="text-center">
        <Spinner animation="border" />
        <p>Mohon tunggu...</p>
      </Container>
    );
  }

  if (error) {
    navigate('/500');
  }

  if (examSessionSerial !== '') {
    if (isSessionAuthorized) {
      navigate(`/exam-session/${examSerial}`);
    }
    return (
      <Container className="text-center">
        <h1 className="my-4 text-center">Ujian {fetchedExam.name}</h1>
        <hr/>
        <QRCodeSVG value={examSessionSerial} size={256}/>
        <hr/>
        <p>Tunjukkan kode QR kepada pengawas untuk memulai ujian.</p>
      </Container>
    );
  }

  const checkIsAuthorized = () => {
    axios.get(`${process.env.REACT_APP_BACKEND_URL}/api/v1/exam-session/${examSerial}/check`)
    .then(response => {
      setIsSessionAuthorized(true);
    })
    .catch(error => {
      setTimeout(checkIsAuthorized, 1000);
    });
  }

  const handleLogin = async (e) => {
    e.preventDefault();

    try {
      const response = await axios.post(`${process.env.REACT_APP_BACKEND_URL}/api/v1/exams/${examSerial}/start`, {
        name: name,
      });
      const token = response.data.data.token;
      localStorage.setItem('examToken', token);
      setExamSessionSerial(response.data.data.session);
      setTimeout(checkIsAuthorized, 1000);
    } catch (error) {
      console.error(error);
      if (error.status === 400) {
        toast.error(`Peserta sudah mengerjakan ujian.`, {
          position: "top-center",
          autoClose: 5000,
          hideProgressBar: false,
          closeOnClick: true,
          pauseOnHover: true,
          draggable: true,
        });
      }
      else if (error.status === 404) {
        toast.error(`Peserta tidak ditemukan. Silakan hubungi pengawas.`, {
          position: "top-center",
          autoClose: 5000,
          hideProgressBar: false,
          closeOnClick: true,
          pauseOnHover: true,
          draggable: true,
        });
      } else {
        toast.error(`Sedang terjadi masalah pada server. Silakan coba beberapa saat lagi.`, {
          position: "top-center",
          autoClose: 5000,
          hideProgressBar: false,
          closeOnClick: true,
          pauseOnHover: true,
          draggable: true,
        });
      }
    }
  }

  return (
    <Container>
      <h1 className="my-4 text-center">Ujian {fetchedExam.name}</h1>
      <hr/>
      <Row className="justify-content-md-center mt-5">
        <Col md={6}>
          <Form onSubmit={handleLogin}>

            <Form.Group controlId="formUsername">
              <Form.Label>Kode</Form.Label>
              <Form.Control
                type="text"
                placeholder="Masukkan kode"
                value={name}
                onChange={(e) => setName(e.target.value)}
                required
              />
            </Form.Group>

            <Button variant="primary" type="submit" className="mt-4 w-100">
              Kerjakan Ujian
            </Button>
          </Form>
        </Col>
      </Row>
    </Container>
  );
}

export default StartExam;