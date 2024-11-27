import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-toastify';
import { Container, Spinner, Button } from 'react-bootstrap';
import QrScanner from 'react-qr-scanner';
import AuthorizeSessionModal from './AuthorizeSessionModal';

const AuthorizeSession = (props) => {
  const { auth } = props;
  const navigate = useNavigate();
  const [scanning, setScanning] = useState(false);
  const [scannedData, setScannedData] = useState("");
  const [showModal, setShowModal] = useState(false);
  const [facingMode, setFacingMode] = useState("environment");

  const handleCloseModal = () => {
    setShowModal(false);
    resumeScanning();
  };

  const handleScan = (data) => {
    if (data) {
      setScannedData(data.text);
      pauseScanning();
      setShowModal(true);
    }
  };

  const pauseScanning = () => {
    setScanning(false);
  }
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

  const toggleCamera = () => {
    setFacingMode((prevMode) => (prevMode === "environment" ? "user" : "environment"));
  };

  return (
    <Container>
      <h1 className="my-4">Izinkan Ujian</h1>
      <hr/>

      <Container className='mt-5 text-center'>
        {scanning ? (
          <Container>
            <Container>
              <Button variant='danger' onClick={pauseScanning}>Hentikan Pemindaian</Button>
            </Container>
            <Container className='mt-5'>
              <QrScanner
                facingMode={facingMode}
                delay={300}
                style={{ width: "100%" }}
                onScan={handleScan}
              />
            </Container>
            <Container>
              <Button onClick={toggleCamera}>
                Ganti kamera
              </Button>
            </Container>
          </Container>
        ) : (
          <Container>
            <Button onClick={resumeScanning}>Pindai Kode QR</Button>
          </Container>
        )}
      </Container>
      <AuthorizeSessionModal
        show={showModal}
        handleClose={handleCloseModal}
        auth={auth}
        examSession={scannedData}
      />
    </Container>
  );
}

export default AuthorizeSession;