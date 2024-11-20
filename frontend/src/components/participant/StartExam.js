import axios from 'axios';
import React, { useEffect, useState } from 'react';
import { Button, Col, Container, Form, Row, Spinner } from 'react-bootstrap';
import { toast } from 'react-toastify';
import { useNavigate, useParams } from 'react-router-dom';

const StartExam = () => {
  const { examSerial } = useParams();
  const navigate = useNavigate();
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(true);
  const [name, setName] = useState('');
  const [password, setPassword] = useState('');
  const [fetchedExam, setFetchedExam] = useState({});

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

  const handleLogin = async (e) => {
    e.preventDefault();

    try {
      const response = await axios.post(`${process.env.REACT_APP_BACKEND_URL}/api/v1/exams/${examSerial}/start`, {
        name: name,
        password: password,
      });
      const token = response.data.data.token;
      localStorage.setItem('examToken', token);
      navigate(`/exam-session/${examSerial}`);
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
              <Form.Label>Nama</Form.Label>
              <Form.Control
                type="text"
                placeholder="Masukkan nama"
                value={name}
                onChange={(e) => setName(e.target.value)}
                required
              />
            </Form.Group>

            <Form.Group controlId="formPassword" className="mt-3">
              <Form.Label>Kata Sandi</Form.Label>
              <Form.Control
                type="password"
                placeholder="Masukkan kata sandi"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                autoComplete='off'
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