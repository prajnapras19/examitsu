import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Card } from 'react-bootstrap';
import { FaHome } from "react-icons/fa";

const BackToHomepageCard = () => {
  const navigate = useNavigate();
  return (
    <Card className="card" onClick={() => navigate('/admin/home')}>
      <Card.Header style={{height: '50%'}}>
        <FaHome style={{height: '100%'}} size={50}></FaHome>
      </Card.Header>
      <Card.Body>
        <Card.Title>
          Menu Utama
        </Card.Title>
        <Card.Text>
          Klik di sini untuk kembali ke menu utama.
        </Card.Text>
      </Card.Body>
    </Card>
  );
}

export default BackToHomepageCard;