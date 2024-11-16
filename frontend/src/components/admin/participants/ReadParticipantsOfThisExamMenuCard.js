import React from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { Card } from 'react-bootstrap';
import { RiFilePaper2Fill } from "react-icons/ri";

const ReadParticipantsOfThisExamMenuCard = () => {
  const navigate = useNavigate();
  const { examSerial } = useParams();
  return (
    <Card className="card" onClick={() => navigate(`/admin/exams/${examSerial}/participants`)}>
      <Card.Header style={{height: '50%'}}>
        <RiFilePaper2Fill style={{height: '100%'}} size={50}></RiFilePaper2Fill>
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