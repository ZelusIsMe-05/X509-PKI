import { useState } from 'react';
import { registerUser, loginUser } from './api/auth';

function App() {
  // State for form inputs and server messages
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');

  // Handle Register button click
  const handleRegister = async () => {
    try {
      const data = await registerUser(username, password);
      setMessage(`Success: ${data.message}`);
    } catch (error: any) {
      setMessage(`Error: ${error.message}`);
    }
  };

  // Handle Login button click
  const handleLogin = async () => {
    try {
      const data = await loginUser(username, password);
      setMessage(`Success: ${data.message}`);
    } catch (error: any) {
      setMessage(`Error: ${error.message}`);
    }
  };

  return (
    <div style={{ padding: '50px', fontFamily: 'sans-serif' }}>
      <h2>API Communication Test: React & Go</h2>
      
      <div style={{ display: 'flex', flexDirection: 'column', width: '300px', gap: '10px' }}>
        <input 
          type="text" 
          placeholder="Enter username" 
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
        <input 
          type="password" 
          placeholder="Enter password" 
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        
        <div style={{ display: 'flex', gap: '10px' }}>
          <button onClick={handleRegister}>Register</button>
          <button onClick={handleLogin}>Login</button>
        </div>
      </div>

      {/* Display server response */}
      <div style={{ marginTop: '20px', padding: '10px', backgroundColor: '#f9f9f9', border: '1px solid #ddd' }}>
        <strong>Server Response:</strong> <br/>
        <span style={{ color: message.startsWith('Error') ? 'red' : 'green', fontWeight: 'bold' }}>
          {message || "Waiting for action..."}
        </span>
      </div>
    </div>
  );
}

export default App;