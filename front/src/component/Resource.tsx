import axios from "axios"
import { useEffect, useState } from 'react';
import AdminPage from "./AdminPage";
import { Routes, Route, useNavigate } from 'react-router-dom';

export const Resource = (props: any) => {
  const { cookies } = props
  const [resource, setResouce] = useState("")
  const [msg, setMsg] = useState("")
  const navigate = useNavigate()

  useEffect(() => {
    axios.get("http://localhost:8000/resource/1234", {
      headers: { Authorization: `Bearer ${cookies.accessToken}`, }
    }
    ).then().catch(err => {
      navigate("/")
    })
  }, [])


  const getCustomer = () => {
    axios.get("http://localhost:8000/resource/1234", {
      headers: { Authorization: `Bearer ${cookies.accessToken}`, }
    }
    ).then(res => {
      setResouce(res.data)
    }).catch(err => {
      if (err.response.data.isAuthorized === false) {
        axios.post("http://localhost:8001/auth/refresh", {
          access_token: cookies.accessToken,
          refresh_token: cookies.refreshToken,
        }).then().catch(err => {
          err.response.data.isAuthorized === false && navigate("/")
        })
      }
    })
  }

  const goToAdminPage = () => {
    axios.get("http://localhost:8000/admin", {
      headers: { Authorization: `Bearer ${cookies.accessToken}`, }
    }
      // 成功した場合は管理者ページに遷移
    ).then(res => {
      navigate("admin")
    }).catch(err => {
      err.response.data.isAuthorized === false && navigate("/resource")
      setMsg(err.response.data.message)
      if (err.response.data.isAuthorized === false) {
        axios.post("http://localhost:8001/auth/refresh", {
          access_token: cookies.accessToken,
          refresh_token: cookies.refreshToken,
        }).then(res => {
          axios.get("http://localhost:8000/admin", {
            headers: { Authorization: `Bearer ${cookies.accessToken}`, }
          })
        }).catch(err => {
          err.response.data.isAuthorized === false && navigate("/")
        })
      }
    })
  }

  return (
    <Routes>
      <Route path="/" element={
        <>
          <h1>ログインに成功しました。</h1>
          <h2>Resourceページ</h2>
          <button onClick={getCustomer}>getCustomer</button>
          <button onClick={goToAdminPage} style={{ marginLeft: "10px" }}>Adminページ</button>
          <p>{resource}</p>
          <p>{msg}</p>
        </>
      }></Route>
      <Route path="admin/*" element={<AdminPage cookies={cookies} />} />
    </Routes>
  )
}

