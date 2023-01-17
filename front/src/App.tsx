import axios from "axios"
import { useState } from "react"
import { Success } from "./component/Success"
import { Routes, Route, useNavigate } from 'react-router-dom';
import { useCookies } from "react-cookie";

function App() {
  const [id, setId] = useState("1234")
  const [password, setPassword] = useState("pass1234")
  const [err, setErr] = useState("")
  const [showErr, setShowErr] = useState(false)
  const navigate = useNavigate()
  const [cookies, setCookie] = useCookies(['accessToken', 'refreshToken']);


  const getAuth = () => {
    axios.post("http://localhost:8001/auth/login", {
      Id: id,
      Password: password,
    }).then(res => {
      setShowErr(false)
      setCookie("accessToken", res.data.access_token, { maxAge: 3600 })
      setCookie("refreshToken", res.data.refresh_token, { maxAge: 3600 })
      navigate("/success")
    }).catch(err => {
      console.log(err)
      console.log(err.response.status)
      console.log(err.response.data)
      setErr(err.response.data)
      setShowErr(true)
    })
  }

  const onChangeId = (e: any) => {
    const id = e.target.value
    setId(id)
  }

  const onChangePassword = (e: any) => {
    const password = e.target.value
    setPassword(password)
  }


  return (
    <Routes>

      <Route path="/" element={
        <div style={{ marginTop: "10px", marginLeft: "5px", width: "300px" }}>
          <table>
            <tbody>
              <tr>
                <td><label htmlFor="id">Id: </label></td>
                <td><input id="id" type="text" onChange={onChangeId} value={id} /></td>
              </tr>
              <tr>
                <td><label htmlFor="password">password: </label></td>
                <td><input id="password" type="text" value={password} onChange={onChangePassword} /></td>
              </tr>
            </tbody>
          </table>
          <button style={{ float: "right" }} onClick={getAuth}>送信</button>
          <p>アクセストークン：{cookies.accessToken}</p>
          <p>リフレッシュトークン：{cookies.refreshToken}</p>
          {showErr && <p>{err}</p>}
        </div>}>
      </Route>
      <Route path="success" element={<Success cookies={cookies} />}></Route>


    </Routes >
  );
}

export default App;
