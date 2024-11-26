import React, { useState, useEffect } from "react";
import axios from 'axios';
import { Container, Modal, Spinner, Button } from "react-bootstrap";

const AuthorizeSessionModal = (props) => {
  const { show, handleClose, auth, examSession } = props;
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (examSession === '') {
      return;
    }
    setLoading(true);
    axios.get(`${process.env.REACT_APP_BACKEND_URL}/api/v1/proctor/participant-sessions/${examSession}/check`, {
      headers: {
        'Authorization': `Bearer ${auth.token}`
      },
    })
    .then(response => {
      setLoading(false);
      setData(response.data.data);
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
    console.log('TODO');
    handleClose();
  }

  console.log('data', data);

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
                  Data: TODO
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
        <Button variant="primary" onClick={handleAuthorize}>
          Izinkan
        </Button>
      </Modal.Footer>
    </Modal>
  )
};

export default AuthorizeSessionModal;