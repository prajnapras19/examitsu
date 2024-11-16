import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Container, Spinner } from 'react-bootstrap';
import LogoutCard from '../auth/LogoutCard';
import ReadExamsMenuCard from '../exam/ReadExamsMenuCard';

const Homepage = (props) => {
  const { auth } = props;
  const navigate = useNavigate();

  useEffect(() => {
    if (auth.loading) {
      return;
    }
    if (!auth.isLoggedIn) {
      navigate('/admin/login');
    }
  }, [auth.loading, auth.isLoggedIn, navigate]);

  if (auth.loading) {
    return (
      <Container className="text-center">
        <Spinner animation="border" />
        <p>Mohon tunggu...</p>
      </Container>
    );  
  }

  return (
    <Container>
      <h1 className="my-4">Menu Utama</h1>
      <hr/>
      <Container className="text-center">
        <p>Selamat datang di Examitsu. Apa yang ingin Anda lakukan hari ini?</p>
      </Container>
      <hr/>
      <Container className="card-grid text-center mt-5">
        <LogoutCard auth={auth}></LogoutCard>
        <ReadExamsMenuCard></ReadExamsMenuCard>
      </Container>
    </Container>
  );
};

export default Homepage;