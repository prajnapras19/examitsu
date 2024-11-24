import React, { useState } from 'react';
import axios from 'axios';
import { Form, Button, Container, Row, Col } from 'react-bootstrap';
import { toast } from 'react-toastify';
import { useNavigate } from 'react-router-dom';

const Login = (props) => {
  const { auth } = props;
  const navigate = useNavigate();
  const [password, setPassword] = useState('');

  const handleLogin = async (e) => {
    e.preventDefault();

    try {
      const response = await axios.post(`${process.env.REACT_APP_BACKEND_URL}/api/v1/proctor/login`, {
        password,
      });
      const token = response.data.data.token;
      localStorage.setItem('authToken', token);
      auth.setLoading(true);
      navigate('/proctor/authorize');
    } catch (error) {
      console.error(error);
      if (error.status === 400) {
        toast.error(`Kata sandi salah. Silakan coba lagi dengan data yang berbeda.`, {
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
  };
  
  if (auth.isLoggedIn) {
    navigate('/proctor/authorize');
  }

  return (
    <Container>
      <Row className="justify-content-md-center mt-5">
        <Col md={6}>
          <h1 className="text-center mb-4">Masuk</h1>
          <Form onSubmit={handleLogin}>

            <Form.Group controlId="formPassword" className="mt-3">
              <Form.Label>Kata Sandi</Form.Label>
              <Form.Control
                type="password"
                placeholder="Masukkan kata sandi (password)"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                autoComplete='off'
                required
              />
            </Form.Group>

            <Button variant="primary" type="submit" className="mt-4 w-100">
              Kirim
            </Button>
          </Form>
        </Col>
      </Row>
    </Container>
  );
};

export default Login;