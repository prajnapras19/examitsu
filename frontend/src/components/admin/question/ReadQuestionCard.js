import React from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { Card } from 'react-bootstrap';
import { FaQuestionCircle } from "react-icons/fa";

const ReadQuestionCard = () => {
  const navigate = useNavigate();
  const { examSerial } = useParams();
  return (
    <Card className="card" onClick={() => navigate(`/admin/exams/${examSerial}/questions`)}>
      <Card.Header style={{height: '50%'}}>
        <FaQuestionCircle style={{height: '100%'}} size={50}></FaQuestionCircle>
      </Card.Header>
      <Card.Body>
        <Card.Title>
          Lihat Daftar Soal
        </Card.Title>
        <Card.Text>
          Klik di sini untuk melihat daftar soal dari ujian ini.
        </Card.Text>
      </Card.Body>
    </Card>
  );
}

export default ReadQuestionCard;