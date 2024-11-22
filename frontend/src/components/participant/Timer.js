import React, { useState, useEffect } from "react";

const Timer = ({ startTime, durationMinutes, onTimesUp }) => {
  const [remainingTime, setRemainingTime] = useState(() => {
    const endTime = new Date(startTime).getTime() + durationMinutes * 60 * 1000;
    return Math.max(0, Math.floor((endTime - Date.now()) / 1000));
  });

  if (remainingTime === 0) {
    onTimesUp();
  }

  useEffect(() => {
    const intervalId = setInterval(() => {
      setRemainingTime((prevTime) => Math.max(0, prevTime - 1));
    }, 1000);

    return () => clearInterval(intervalId); // Cleanup on unmount
  }, []);

  const formatTime = (seconds) => {
    const minutes = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${minutes.toString().padStart(2, "0")}:${secs.toString().padStart(2, "0")}`;
  };

  return (
    <div>
      <p>{formatTime(remainingTime)}</p>
    </div>
  );
};

export default Timer;
