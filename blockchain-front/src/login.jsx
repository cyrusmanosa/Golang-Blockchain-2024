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
    const [message, setMessage] = useState('')
    const [file, setFile] = useState(null);
    const [company_name, setCompanyName] = useState('')
    const [hash, setHash] = useState('');


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

    const handleChange = useCallback((e) => {
        setHash(e.target.value);
    });


    const handleSendData = async () => {
        const formData = new FormData();
        const formFile = new FormData();
        formData.append('name', name);
        formData.append('email', email);
        formData.append('company_name', company_name);
        formData.append('message', message);
        formData.append('hash', hash);

        console.log(hash);

        await axios.post("http://localhost:8080/take", formData, {
            headers: {
                'Content-Type': 'application/json'
            }
        });

        if (file != null) {
            formFile.append('file', file);
            await axios.post("http://localhost:8080/Upload", formFile, {
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

            <label htmlFor="hash-algorithm">Select Hash : </label>
            <select
                id="hash-algorithm"
                value={hash}
                onChange={handleChange}
            >
                <option value="">--Select--</option>
                <option value="sha256">SHA-256</option>
                {/* <option value="argon2">Argon2</option> */}
                <option value="blake2b">Blake2b</option>
                {/* <option value="blake2s">Blake2s</option> */}
                <option value="blake3">Blake3</option>
                <option value="gost-r">GOST R</option>
                {/* <option value="ripemd160">RIPEMD-160</option> */}
                <option value="keccak">Keccak</option>
                <option value="skein">Skein</option>
                <option value="skein">MurmurHash3</option>
                <option value="skein">FarmHash</option>
                <option value="skein">xxHash</option>
                <option value="skein">HighwayHash</option>
                {/* <option value="whirlpool">Whirlpool</option> */}
            </select>

            <br />
            <button onClick={handleSendData}>提交</button>
        </Box>
    );
}
