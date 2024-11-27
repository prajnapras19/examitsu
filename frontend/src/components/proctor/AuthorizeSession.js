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
  const [state, setState] = useState({ cameraId: undefined, delay: 500, devices: [], loading: false })

  navigator.mediaDevices.enumerateDevices()
  .then((devices) => {
    const videoSelect = []
    devices.forEach((device) => {
      if (device.kind === 'videoinput') {
        videoSelect.push(device)
      }
    })
    return videoSelect
  })
  .then((devices) => {
    setState({
      cameraId: devices[0].deviceId,
      devices,
      loading: false,
    })
  })
  .catch((error) => {
    toast.error(`Terjadi kesalahan.`, {
      position: "top-center",
      autoClose: 5000,
      hideProgressBar: false,
      closeOnClick: true,
      pauseOnHover: true,
      draggable: true,
    });
  });

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
    setState({
      loading: true,
    })
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


  const { loading, cameraId, devices } = state;
  return (
    <Container>
      <h1 className="my-4">Izinkan Ujian</h1>
      <hr/>

      <Container className='mt-5 text-center'>
        {scanning ? (
          <Container>
            {state.loading
            ? (
              <Container className="text-center">
                <Spinner animation="border" />
                <p>Mohon tunggu...</p>
              </Container>
            )
            : (
                <>
                  <Container>
                    <Button variant='danger' onClick={pauseScanning}>Hentikan Pemindaian</Button>
                  </Container>
                  <Container className="mt-5 flex">
                    <Container>
                      Kamera:
                    </Container>
                    <Container>
                      <select
                        onChange={e => {
                          const value = e.target.value
                          this.setState({ cameraId: undefined }, () => {
                            this.setState({ cameraId: value })
                          })
                        }}
                      >
                        {devices.map((deviceInfo, index) => (
                          <React.Fragment key={deviceInfo.deviceId}><option value={deviceInfo.deviceId}>{deviceInfo.label || `camera ${index}`}</option></React.Fragment>
                        ))}
                      </select>
                    </Container>
                  </Container>
                  <Container className='mt-5'>
                    <QrScanner
                      delay={300}
                      style={{ width: "100%" }}
                      onScan={handleScan}
                    />
                  </Container>
                </>
              )
            }
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