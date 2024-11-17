import { Route, Routes } from 'react-router-dom';
import NotFoundPage from './components/etc/404';
import ExamSession from './components/participant/ExamSession';

function ParticipantRoutes() {
  document.addEventListener('contextmenu', (e) => {
    e.preventDefault();
  });
  return (
    <Routes>
      <Route path="/:examSerial" element={<ExamSession/>}/>
      <Route path="*" element={<NotFoundPage/>}/>
    </Routes>
  );
}

export default ParticipantRoutes;