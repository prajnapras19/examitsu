import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-toastify';
import { Container, Spinner, Button } from 'react-bootstrap';

const AuthorizeSession = (props) => {
  const { auth } = props;
  const navigate = useNavigate();
  const [error, setError] = useState(null);
  
  useEffect(() => {
    if (auth.loading) {
      return;
    }
    if (!auth.isLoggedIn) {
      navigate('/proctor/login');
    }
  }, [auth.loading, auth.isLoggedIn]);

  if (auth.loading) {
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

  return (
    <Container>
      <h1 className="my-4">Izinkan Ujian</h1>
      <hr/>


      <p>TODO</p>
    </Container>
  );
}

export default AuthorizeSession;