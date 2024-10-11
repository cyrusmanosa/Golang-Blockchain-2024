import { useParams } from 'react-router-dom';
import axios from 'axios';

export default function Check() {
    const { name } = useParams();
    const handleSendData = () => {
        axios.put(`http://localhost:8080/Check/${name}`,);
    }
    handleSendData();
}