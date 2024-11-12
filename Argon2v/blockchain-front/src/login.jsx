import Box from '@mui/material/Box';
import InputAdornment from '@mui/material/InputAdornment';
import TextField from '@mui/material/TextField';
import AccountCircle from '@mui/icons-material/AccountCircle';
import CloudUploadIcon from '@mui/icons-material/CloudUpload';
import { useState, useCallback } from 'react';
import axios from 'axios';

export default function InputWithIcon() {
    const [name, setName] = useState('')
    const [email, setEmail] = useState('')
    const [company_name, setCompanyName] = useState('')
    const [message, setMessage] = useState('')
    const [file, setFile] = useState(null);

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

    const handleFileChange = useCallback((e) => {
        setFile(e.target.files[0]);
    }, []);

    const handleSendData = async () => {
        const formData = new FormData();
        const formFile = new FormData();
        formData.append('name', name);
        formData.append('email', email);
        formData.append('company_name', company_name);
        formData.append('message', message);

        await axios.post("http://localhost:8081/take", formData, {
            headers: {
                'Content-Type': 'application/json'
            }
        });

        if (file != null) {
            formFile.append('file', file);
            await axios.post("http://localhost:8081/Upload", formFile, {
                headers: {
                    'Content-Type': 'application/pdf'
                }
            });
        }
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

            {/* 文件上传 */}
            <TextField
                id="File"
                type="file"
                InputProps={{
                    startAdornment: (
                        <InputAdornment position="start">
                            <CloudUploadIcon />
                        </InputAdornment>
                    ),
                }}
                variant="standard"
                onChange={handleFileChange}
            />

            <br />
            <button onClick={handleSendData}>提交</button>
        </Box>
    );
}
