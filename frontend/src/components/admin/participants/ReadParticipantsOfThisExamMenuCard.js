import React from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { Card } from 'react-bootstrap';
import { BsPeopleFill } from "react-icons/bs";

const ReadParticipantsOfThisExamMenuCard = () => {
  const navigate = useNavigate();
  const { examSerial } = useParams();
  return (
    <Card className="card" onClick={() => navigate(`/admin/exams/${examSerial}/participants`)}>
      <Card.Header style={{height: '50%'}}>
        <BsPeopleFill style={{height: '100%'}} size={50}></BsPeopleFill>
      </Card.Header>
      <Card.Body>
        <Card.Title>
          Lihat Daftar Peserta
        </Card.Title>
        <Card.Text>
          Klik di sini untuk melihat daftar peserta dari ujian ini.
        </Card.Text>
      </Card.Body>
    </Card>
  );
}

export default ReadParticipantsOfThisExamMenuCard;