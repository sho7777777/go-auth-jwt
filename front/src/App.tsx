import axios from "axios"
import { useState } from "react"
import { Resource } from "./component/Resource"
import { Routes, Route, useNavigate } from 'react-router-dom';
import { useCookies } from "react-cookie";

function App() {
  // 認証情報
  const [id, setId] = useState("1234")
  const [password, setPassword] = useState("pass1234")
  // エラー
  const [err, setErr] = useState("")
  const [showErr, setShowErr] = useState(false)
  // クッキー管理
  const [cookies, setCookies] = useCookies(['accessToken', 'refreshToken']);
  // ルーティング
  const navigate = useNavigate()

  const getAuth = () => {
    axios.post("http://localhost:8001/auth/login", {
      Id: id,
      Password: password,
      // 成功した場合、取得したトークンをクッキーをセットし、resourceページに遷移する、
    }).then(res => {
      setShowErr(false)
      setCookies("accessToken", res.data.access_token, { maxAge: 3600 })
      setCookies("refreshToken", res.data.refresh_token, { maxAge: 3600 })
      navigate("/resource")
      // 失敗した場合はエラーメッセージを表示
    }).catch(err => {
      if (err.code == "ERR_NETWORK") {
        setErr("ネットワークにエラーが発生しました。")
      } else {
        setErr(err.response.data)
      }
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
      <Route path="resource/*" element={<Resource cookies={cookies} setCookies={setCookies} />}></Route>
    </Routes >
  );
}

export default App;
