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

function App() {
  return (
    <div className="App my-4">
      <Router>
        <ScrollToTop />
        <ToastContainer />
        <Routes>
          <Route path="/admin/*" element={<AdminRoutes/>}/>
          <Route path="/public/exams/:examSerial" element={<StartExam/>} />
          <Route path="/public/*" element={<ParticipantRoutes/>}/>
          <Route path="/500" element={<InternalServerErrorPage/>} />
          <Route path="*" element={<NotFoundPage/>}/>
        </Routes>
      </Router>
    </div>
  );
}

export default App;
