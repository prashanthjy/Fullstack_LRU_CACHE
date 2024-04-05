import React, { useState } from 'react';
import axios from 'axios';

const App = () => {
    const [key, setKey] = useState('');
    const [value, setValue] = useState('');
    const [ttl, setTtl] = useState('');
    const [getResponse, setGetResponse] = useState('');
    const [postResponse, setPostResponse] = useState('');

    const handleGetCache = async () => {
        try {
            const response = await axios.get(`http://localhost:8080/cache/${key}`);
            setGetResponse(response.data.value)
        } catch (error) {
            setGetResponse(error.response.data);
        }
    };

    const handleSetCache = async () => {
        try {
            await axios.post('http://localhost:8080/cache', {
                key,
                value,
                ttl: parseInt(ttl)
            });
            setPostResponse(`'Value set successfully for key', ${key}`)
        } catch (error) {
            setPostResponse(error.response.data);
        }
    };

    return (
        <div>
            <h1>LRU Cache Frontend</h1>
            <div>
                <h2>Get Value from Cache</h2>
                <input type="text" value={key} onChange={(e) => setKey(e.target.value)} />
                <button onClick={handleGetCache}>Get Value</button>
                <p>Response: {getResponse}</p>
            </div>
            <div>
                <h2>Set Value in Cache</h2>
                <input type="text" placeholder="Key" value={key} onChange={(e) => setKey(e.target.value)} />
                <input type="text" placeholder="Value" value={value} onChange={(e) => setValue(e.target.value)} />
                <input type="number" placeholder="TTL" value={ttl} onChange={(e) => setTtl(e.target.value)} />
                <button onClick={handleSetCache}>Set Value</button>
                <p>Response: {postResponse}</p>
            </div>
        </div>
    );
};

export default App;