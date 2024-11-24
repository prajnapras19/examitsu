import React, { useEffect } from 'react';
import axios from 'axios';
import { QRCodeSVG } from 'qrcode.react';
import { Container } from 'react-bootstrap';

const ExamSessionQRCode = ({ fetchedExam, examSessionSerial, setIsSessionAuthorized, examSerial }) => {
  const checkIsAuthorized = () => {
    const token = localStorage.getItem('examToken');
    axios.get(`${process.env.REACT_APP_BACKEND_URL}/api/v1/exam-session/${examSerial}/check`, {
      headers: {
        'Authorization': `Bearer ${token}`
      },
    })
    .then(response => {
      setIsSessionAuthorized(true);
    })
    .catch(error => {
      // pass
    });
  }

  useEffect(() => {
    const intervalId = setInterval(checkIsAuthorized, 1000);
    return () => {
      clearInterval(intervalId);
    }
  }, []);

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

export default ExamSessionQRCode;