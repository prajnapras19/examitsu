import React, { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import InternalServerErrorPage from '../etc/500';
import { Button, Container, Spinner, Form, Row, Col } from 'react-bootstrap';
import axios from 'axios';
import QuestionListSidebar from './QuestionListSidebar';
import EditorJsHTML from 'editorjs-html';
import { toast } from 'react-toastify';
import FinishPage from './FinishPage';
import SubmitConfirmationModal from './SubmitConfirmationModal';
import Timer from './Timer';
import { shuffleArray } from '../../utils/random';

const ExamSession = () => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const { examSerial } = useParams();
  const navigate = useNavigate();
  const [questionIDList, setQuestionIDList] = useState([]);
  const [questionIDSet, setQuestionIDSet] = useState(new Set());
  const [currentQuestionNumber, setCurrentQuestionNumber] = useState(1);
  const [currentQuestion, setCurrentQuestion] = useState(null);
  const edjsParser = EditorJsHTML();
  const [disableChooseOption, setDisableChooseOption] = useState(false);
  const [disableChangeQuestion, setDisableChangeQuestion] = useState(false);
  const [startTime, setStartTime] = useState(null);
  const [duration, setDuration] = useState(null);

  const [loadingSubmit, setLoadingSubmit] = useState(false);
  const [isSubmitted, setIsSubmitted] = useState(false);
  const [showSubmitModal, setShowSubmitModal] = useState(false);
  const handleShowSubmitModal = () => {
    setShowSubmitModal(true);
  }
  const handleCloseSubmitModal = () => {
    setShowSubmitModal(false);
  }

  const handleSubmit = () => {
    handleCloseSubmitModal();
    setLoadingSubmit(true);
    const token = localStorage.getItem('examToken');

    if (!token) {
      navigate('/404');
    }
    axios.post(`${process.env.REACT_APP_BACKEND_URL}/api/v1/exam-session/${examSerial}/submit`,
      {},
      { 
        headers: {
        'Authorization': `Bearer ${token}`
        },
      }
    )
    .then(response => {
      setIsSubmitted(true);
    })
    .catch(err => {
      toast.error(`Gagal mengumpulkan ujian. Silakan coba beberapa saat lagi.`, {
        position: "top-center",
        autoClose: 5000,
        hideProgressBar: false,
        closeOnClick: true,
        pauseOnHover: true,
        draggable: true,
      });
    }).finally(() => {
      setLoadingSubmit(false);
    });
  }

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
      const newSet = new Set();
      for (let i = 0; i < response.data.data.questions_id_list.length; i++) {
        newSet.add(response.data.data.questions_id_list[i].id);
      }

      let shuffledQuestions = questionIDList;
      if (newSet.symmetricDifference(questionIDSet).size > 0) {
        setQuestionIDSet(newSet);
        shuffledQuestions = shuffleArray(response.data.data.questions_id_list);
        setQuestionIDList(shuffledQuestions);
      }
      
      setStartTime(response.data.data.start_time);
      setDuration(response.data.data.duration);
      setDisableChooseOption(false);
      
      if (currentQuestionNumber > shuffledQuestions.length) {
        setCurrentQuestion(null);
      } else {
        axios.get(`${process.env.REACT_APP_BACKEND_URL}/api/v1/exam-session/${examSerial}/questions/${shuffledQuestions[currentQuestionNumber-1].id}`, {
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

  if (loading || loadingSubmit) {
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

  if (isSubmitted) {
    return (
      <FinishPage></FinishPage>
    )
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
    setDisableChangeQuestion(true);
    setTimeout(() => {
      setDisableChangeQuestion(false);
    }, 1000);
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
  
  const handleClickOption = async (optionId) => {
    setDisableChooseOption(true);
    try {
      const token = localStorage.getItem('examToken');

      if (!token) {
        navigate('/404');
      }

      await axios.post(`${process.env.REACT_APP_BACKEND_URL}/api/v1/exam-session/${examSerial}/questions/${currentQuestion.question.id}`,
        {
          mcq_option_id: optionId,
        },
        {
          headers: {
            'Authorization': `Bearer ${token}`
          },
        }
      );
      setTimeout(() => {
        toast.success(`Jawaban untuk nomor ${currentQuestionNumber} berhasil disimpan!`, {
          position: "top-center",
          autoClose: 3000,
          hideProgressBar: false,
          closeOnClick: true,
          pauseOnHover: true,
          draggable: true,
        });
      }, 1000);
    } catch (err) {
      if (err.status < 500) {
        navigate('/404');
      } else {
        setTimeout(() => {
          toast.error(`Gagal menyimpan jawaban untuk nomor ${currentQuestionNumber}, silakan coba beberapa saat lagi.`, {
            position: "top-center",
            autoClose: 5000,
            hideProgressBar: false,
            closeOnClick: true,
            pauseOnHover: true,
            draggable: true,
          });
        }, 1000);
      }
    } finally {
      setTimeout(() => {
        setDisableChooseOption(false)
      }, 500);
    }
  }

  return (
    <>
      <hr className='mt-5'/>
        <Container className="prevent-select">
          <QuestionListSidebar
            questionIDList={questionIDList}
            handleChooseQuestion={handleChooseQuestion}
            handleShowSubmitModal={handleShowSubmitModal}
            disabled={disableChangeQuestion}
          />
        </Container>
      <hr/>
      <Container className="prevent-select">
        <p>
          <b>Waktu tersisa:</b> <Timer startTime={startTime} durationMinutes={duration} onTimesUp={() => setIsSubmitted(true)}></Timer>
        </p>
      </Container>
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
                      disabled={disableChooseOption}
                    />
                  ))}
                </Form>
              )
              : (
                <></>
              )
            }
          </Container>
          <SubmitConfirmationModal
            show={showSubmitModal}
            handleClose={handleCloseSubmitModal}
            handleSubmit={handleSubmit}
          />
          <Container className="d-flex mt-3 prevent-select">
            <Button
              className="me-3"
              variant="primary"
              onClick={() => handleChooseQuestion(currentQuestionNumber - 1)}
              disabled={currentQuestionNumber === 1 || disableChangeQuestion}
            >
              &lt;&lt; Soal sebelumnya
            </Button>
            <Button
              variant="primary"
              onClick={() => handleChooseQuestion(currentQuestionNumber + 1)}
              disabled={currentQuestionNumber === questionIDList.length || disableChangeQuestion}
            >
              Soal selanjutnya &gt;&gt;
            </Button>
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