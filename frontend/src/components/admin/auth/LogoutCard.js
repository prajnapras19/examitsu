import React from 'react';
import { Card } from 'react-bootstrap';
import { GrLogout } from "react-icons/gr";

const LogoutCard = (props) => {
  const { auth } = props;

  const logout = () => {
    localStorage.removeItem('authToken');
    auth.setLoading(true);
  }

  return (
    <Card className="card" onClick={logout}>
        <Card.Header style={{height: '50%'}}>
          <GrLogout style={{height: '100%'}} size={50}></GrLogout>
        </Card.Header>
        <Card.Body>
          <Card.Title>
            Keluar Dari Aplikasi
          </Card.Title>
          <Card.Text>
            Klik di sini untuk keluar dari aplikasi.
          </Card.Text>
        </Card.Body>
    </Card>
  );
}

export default LogoutCard;