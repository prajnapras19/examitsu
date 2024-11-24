import './App.css';
import 'bootstrap/dist/css/bootstrap.min.css';
import 'react-toastify/dist/ReactToastify.css';
import { ToastContainer } from 'react-toastify';
import NotFoundPage from './components/etc/404';
import ScrollToTop from './components/etc/ScrollToTop';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import AdminRoutes from './AdminRoutes';
import ParticipantRoutes from './ParticipantRoutes';
import StartExam from './components/participant/StartExam';
import InternalServerErrorPage from './components/etc/500';
import GetAllOpenedExams from './components/participant/GetAllOpenedExams';
import ProctorRoutes from './ProctorRoutes';

function App() {
  return (
    <div className="App my-4">
      <Router>
        <ScrollToTop />
        <ToastContainer />
        <Routes>
          <Route path="/admin/*" element={<AdminRoutes/>}/>
          <Route path="/" element={<GetAllOpenedExams/>} />
          <Route path="/exam/:examSerial" element={<StartExam/>} />
          <Route path="/exam-session/*" element={<ParticipantRoutes/>}/>
          <Route path="/proctor/*" element={<ProctorRoutes/>}/>
          <Route path="/500" element={<InternalServerErrorPage/>} />
          <Route path="*" element={<NotFoundPage/>}/>
        </Routes>
      </Router>
    </div>
  );
}

export default App;
