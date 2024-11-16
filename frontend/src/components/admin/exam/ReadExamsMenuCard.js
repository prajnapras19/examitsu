import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Card } from 'react-bootstrap';
import { RiFilePaper2Fill } from "react-icons/ri";

const ReadExamsMenuCard = () => {
  const navigate = useNavigate();
  return (
    <Card className="card" onClick={() => navigate('/admin/exams')}>
      <Card.Header style={{height: '50%'}}>
        <RiFilePaper2Fill style={{height: '100%'}} size={50}></RiFilePaper2Fill>
      </Card.Header>
      <Card.Body>
        <Card.Title>
          Lihat Ujian
        </Card.Title>
        <Card.Text>
          Klik di sini untuk melihat ujian-ujian yang sudah didaftarkan ke dalam sistem.
        </Card.Text>
      </Card.Body>
    </Card>
  );
}

export default ReadExamsMenuCard;