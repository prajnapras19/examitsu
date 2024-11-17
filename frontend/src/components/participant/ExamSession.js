import React, { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import InternalServerErrorPage from '../etc/500';
import { Container, Spinner } from 'react-bootstrap';
import axios from 'axios';
import QuestionListSidebar from './QuestionListSidebar';

const ExamSession = () => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const { examSerial } = useParams();
  const navigate = useNavigate();
  const [questionIDList, setQuestionIDList] = useState([]);
  const [currentQuestionNumber, setCurrentQuestionNumber] = useState(1);

  useEffect(() => {
    const token = localStorage.getItem('examToken');

    if (!token) {
      navigate('/404');
    }

    axios.get(`${process.env.REACT_APP_BACKEND_URL}/api/v1/exam-session/${examSerial}/questions`, {
      headers: {
        'Authorization': `Bearer ${token}`
      },
    })
    .then(response => { 
      setQuestionIDList(response.data.data);
      setLoading(false);
    })
    .catch(error => {
      if (error.status < 500) {
        navigate('/404');
      }
      else {
        setError(error.message);
      }
      setLoading(false);
    });
  }, [currentQuestionNumber]);

  if (loading) {
    return (
      <Container className="text-center">
        <Spinner animation="border" />
        <p>Mohon tunggu...</p>
      </Container>
    );
  }

  if (error) {
    return <InternalServerErrorPage></InternalServerErrorPage>
  }

  if (questionIDList.length == 0) {
    return (
      <Container className="text-center mt-5 prevent-select">
        <p>
          <i>Tidak ada soal tersedia.</i>
        </p>
      </Container>
    )
  }

  const handleChooseQuestion = (id) => {
    setCurrentQuestionNumber(id);
  }

  return (
    <>
      <hr/>
      <QuestionListSidebar
        questionIDList={questionIDList}
        handleChooseQuestion={handleChooseQuestion}
      />
      <hr/>
      <Container className="text-center mt-5 prevent-select">  
      </Container>
    </>
  );
}

export default ExamSession;