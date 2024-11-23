import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import { Container, Spinner, Table, Button } from 'react-bootstrap';

const GetAllOpenedExams = () => {
  const navigate = useNavigate();

  const [currentPage, setCurrentPage] = useState(1);
  const [data, setData] = useState([]);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get(`${process.env.REACT_APP_BACKEND_URL}/api/v1/exams?page=${currentPage}`);
        setData(response.data.data);
        setLoading(false);
      } catch (err) {
        console.error("Error fetching data", err);
        setError(err);
        setLoading(false);
      }
    };

    fetchData();
  }, [currentPage]);

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

  const handleNextPage = () => {
    setCurrentPage(currentPage + 1);
  };

  const handlePreviousPage = () => {
    if (currentPage > 1) {
      setCurrentPage(currentPage - 1);
    }
  };

  return (
    <Container>
      <h1 className="my-4">Daftar Ujian</h1>
      <hr/>

      {data.length === 0 ? (
        <Container className="text-center mt-5">
          {currentPage === 1 ? (
            <i>Tidak ada data ditemukan.</i>  
          ): (
            <i>Tidak ada data ditemukan pada halaman ini. Silakan coba halaman sebelumnya.</i>
          )}
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
      <Container className="d-flex mt-3">
        <Button variant="primary" onClick={handlePreviousPage} disabled={currentPage === 1} className="me-3">
          Halaman sebelumnya
        </Button>
        <Button variant="primary" onClick={handleNextPage} disabled={data.length === 0}>
          Halaman berikutnya
        </Button>
      </Container>
    </Container>
  );
}

export default GetAllOpenedExams;