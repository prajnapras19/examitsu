import axios from 'axios';
import React, { useEffect, useState } from 'react';
import { Container, Spinner } from 'react-bootstrap';
import { useNavigate, useParams } from 'react-router-dom';

const StartExam = () => {
  const { examSerial } = useParams();
  const navigate = useNavigate();
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(true);
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

  return (
    <Container>
      <h1 className="my-4 text-center">Ujian {fetchedExam.name}</h1>
      <hr/>
    </Container>
  );
}

export default StartExam;