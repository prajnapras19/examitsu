import React, { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import InternalServerErrorPage from '../etc/500';
import { Container, Spinner, Form } from 'react-bootstrap';
import axios from 'axios';
import QuestionListSidebar from './QuestionListSidebar';
import EditorJsHTML from 'editorjs-html';

const ExamSession = () => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const { examSerial } = useParams();
  const navigate = useNavigate();
  const [questionIDList, setQuestionIDList] = useState([]);
  const [currentQuestionNumber, setCurrentQuestionNumber] = useState(1);
  const [currentQuestion, setCurrentQuestion] = useState(null);
  const edjsParser = EditorJsHTML();
  const [selectedOption, setSelectedOption] = useState(0);

  useEffect(() => {
    if (!loading) {
      return;
    }
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
      setSelectedOption(0);
      
      if (currentQuestionNumber > response.data.data.length) {
        setCurrentQuestion(null);
      } else {
        axios.get(`${process.env.REACT_APP_BACKEND_URL}/api/v1/exam-session/${examSerial}/questions/${response.data.data[currentQuestionNumber-1].id}`, {
          headers: {
            'Authorization': `Bearer ${token}`
          },
        })
        .then(response => { 
          setCurrentQuestion(response.data.data);
        }).catch(error => {
          setCurrentQuestion(null);
          setLoading(false);
        });
      }
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
  }, [loading]);

  useEffect(() => {
    setCurrentQuestion(null);
    setLoading(true);
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

  const handleChooseQuestion = (i) => {
    setCurrentQuestionNumber(i);
  }

  let parsedHTML = undefined;
  if (currentQuestion) {
    try {
      parsedHTML = edjsParser.parse(JSON.parse(currentQuestion.question.data)).join('');
    } catch (err) {
      parsedHTML = `<div></div>`;
    }
  }
  
  const handleClickOption = (optionId) => {
    console.log('optionId', optionId);
  }

  return (
    <>
      <hr/>
      <QuestionListSidebar
        questionIDList={questionIDList}
        handleChooseQuestion={handleChooseQuestion}
      />
      <hr/>
      <h3 className="text-center prevent-select">Soal {currentQuestionNumber}</h3>
      <hr/>
      <Container className="mt-5 prevent-select">    
        {currentQuestion
        ? (
          <div dangerouslySetInnerHTML={{ __html: parsedHTML }} />
        )
        : (
          <p>
            <i>Soal tidak ditemukan.</i>
          </p>
        )}
      </Container>
      {currentQuestion
      ? (
        <>
          <hr/>
          <Container className="mt-3 prevent-select">
            <h6>Pilihan Jawaban:</h6>
          </Container>
          <hr/>
          <Container className="mt-3 prevent-select">
            {currentQuestion.options
              ? (
                <Form className="ms-3 px-3" style={{fontSize: 20}}>
                  {currentQuestion.options.map((data) => (
                    <Form.Check
                      size='lg'
                      type='radio'
                      name='option'
                      label={data.description}
                      onClick={() => handleClickOption(data.id)}
                      defaultChecked={data.id === currentQuestion.answer}
                    />
                  ))}
                </Form>
              )
              : (
                <></>
              )
            }
          </Container>
        </>
      )
      : (
        <></>
      )}
      
    </>
  );
}

export default ExamSession;