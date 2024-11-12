import { useParams } from 'react-router-dom';
import axios from 'axios';

export default function Check() {
    const { name } = useParams();
    axios.post(`http://localhost:8081/Check/${name}`,);
}