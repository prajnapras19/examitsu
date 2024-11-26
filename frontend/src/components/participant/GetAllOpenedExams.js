import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import { Container, Spinner, Table, Button } from 'react-bootstrap';

const GetAllOpenedExams = () => {
  const navigate = useNavigate();

  const [data, setData] = useState([]);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get(`${process.env.REACT_APP_BACKEND_URL}/api/v1/exams`);
        setData(response.data.data);
        setLoading(false);
      } catch (err) {
        console.error("Error fetching data", err);
        setError(err);
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  if (loading) {
    return (
      <Container className="text-center">
        <Spinner animation="border" />
        <p>Mohon tunggu...</p>
      </Container>
    );
  }

  if (error) {
    navigate('/500');
  }

  return (
    <Container>
      <h1 className="my-4">Daftar Ujian</h1>
      <hr/>

      {data.length === 0 ? (
        <Container className="text-center mt-5">
          <i>Tidak ada data ditemukan.</i>
        </Container>
      ) : (
        <Table striped bordered hover className="text-center mt-5">
          <thead>
            <tr>
              <th>Nama</th>
              <th>Durasi Pengerjaan (menit)</th>
              <th colSpan="5">Aksi</th>
            </tr>
          </thead>
          <tbody>
            {data.map((exam) => (
              <tr key={exam.serial}>
                <td className="p-3">{exam.name}</td>
                <td className="p-3">{exam.allowed_duration_minutes}</td>
                <td>
                  <Button variant="primary" className="me-3" onClick={() => navigate(`/exam/${exam.serial}`)}>Kerjakan Ujian</Button>
                </td>
              </tr>
            ))}
          </tbody>
        </Table>
      )}
    </Container>
  );
}

export default GetAllOpenedExams;