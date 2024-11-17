import React from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { Card } from 'react-bootstrap';
import { FaPlus } from "react-icons/fa";

const AddParticipantsCard = () => {
  const navigate = useNavigate();
  const { examSerial } = useParams();
  return (
    <Card className="card" onClick={() => navigate(`/admin/exams/${examSerial}/participants/new`)}>
      <Card.Header style={{height: '50%'}}>
        <FaPlus style={{height: '100%'}} size={50}></FaPlus>
      </Card.Header>
      <Card.Body>
        <Card.Title>
          Tambah Peserta
        </Card.Title>
        <Card.Text>
          Klik di sini untuk menambah peserta untuk ujian ini ke dalam sistem.
        </Card.Text>
      </Card.Body>
    </Card>
  );
}

export default AddParticipantsCard;