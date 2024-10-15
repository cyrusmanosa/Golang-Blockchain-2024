import Box from '@mui/material/Box';
import InputAdornment from '@mui/material/InputAdornment';
import TextField from '@mui/material/TextField';
import AccountCircle from '@mui/icons-material/AccountCircle';
import { useState, useCallback } from 'react';
import axios from 'axios';

export default function InputWithIcon() {
    const [name, setName] = useState('')
    const [email, setEmail] = useState('')
    const [company_name, setCompanyName] = useState('')
    const [message, setMessage] = useState('')

    const handleName = useCallback((e) => {
        setName(e.target.value);
    }, []);

    const handleEmail = useCallback((e) => {
        setEmail(e.target.value);
    }, []);

    const handleCompanyName = useCallback((e) => {
        setCompanyName(e.target.value);
    }, []);

    const handleMessage = useCallback((e) => {
        setMessage(e.target.value);
    }, []);


    const handleSendData = async () => {
        const data = {
            name,
            email,
            company_name,
            message
        };

        await axios.post("http://localhost:8080/take", data, {
            headers: {
                'Content-Type': 'application/json'
            }
        });
    }

    return (
        <Box sx={{ '& > :not(style)': { m: 1 } }}>
            {/* Name */}
            <TextField
                id="Name"
                label="Name"
                InputProps={{
                    startAdornment: (
                        <InputAdornment position="start">
                            <AccountCircle />
                        </InputAdornment>
                    ),
                }}
                variant="standard"
                value={name}
                onChange={handleName}
            />
            {/* Email */}
            <TextField
                id="Email"
                label="Email"
                InputProps={{
                    startAdornment: (
                        <InputAdornment position="start">
                            <AccountCircle />
                        </InputAdornment>
                    ),
                }}
                variant="standard"
                value={email}
                onChange={handleEmail}
            />
            {/* CompanyName */}
            <TextField
                id="CompanyName"
                label="CompanyName"
                InputProps={{
                    startAdornment: (
                        <InputAdornment position="start">
                            <AccountCircle />
                        </InputAdornment>
                    ),
                }}
                variant="standard"
                value={company_name}
                onChange={handleCompanyName}
            />
            {/* Message */}
            <TextField
                id="Message"
                label="Message"
                InputProps={{
                    startAdornment: (
                        <InputAdornment position="start">
                            <AccountCircle />
                        </InputAdornment>
                    ),
                }}
                variant="standard"
                value={message}
                onChange={handleMessage}
            />

            <br />
            <button onClick={handleSendData}>Submit</button>
        </Box>
    );
}
