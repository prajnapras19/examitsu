import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-toastify';
import { Container, Spinner, Button } from 'react-bootstrap';
import QrScanner from 'react-qr-scanner';

const AuthorizeSession = (props) => {
  const { auth } = props;
  const navigate = useNavigate();
  const [scanning, setScanning] = useState(false);
  const [scannedData, setScannedData] = useState("");

  const handleScan = (data) => {
    if (data) {
      setScannedData(data.text); // Save the scanned data
      setScanning(false);       // Pause scanning
      // Perform your desired action here
      console.log("Scanned Data:", data.text);
    }
  };

  const resumeScanning = () => {
    setScanning(true);
    setScannedData("");
  };
  
  useEffect(() => {
    if (auth.loading) {
      return;
    }
    if (!auth.isLoggedIn) {
      navigate('/proctor/login');
    }
  }, [auth.loading, auth.isLoggedIn]);

  if (auth.loading) {
    return (
      <Container className="text-center">
        <Spinner animation="border" />
        <p>Mohon tunggu...</p>
      </Container>
    );
  }

  return (
    <Container>
      <h1 className="my-4">Izinkan Ujian</h1>
      <hr/>

      <Container className='mt-5 text-center'>
        {scanning ? (
          <QrScanner
            delay={300}
            style={{ width: "30%" }}
            onScan={handleScan}
          />
        ) : (
          <div>
            <h2>Scanned Data: {scannedData}</h2>
            <button onClick={resumeScanning}>Scan Again</button>
          </div>
        )}
      </Container>
    </Container>
  );
}

export default AuthorizeSession;