import axios from "axios"
import { useNavigate } from 'react-router-dom';
import { useEffect, useState } from 'react';

export const Success = (props: any) => {
  const { cookies } = props
  const [msg, setMsg] = useState("")
  const navigate = useNavigate()

  useEffect(() => {
    if (cookies.accessToken === undefined || cookies.accessToken === "") {
      navigate("/")
    }
  }, [])


  const getCustomer = () => {
    axios.get("http://localhost:8000/resource/1234", {
      headers: { Authorization: `Bearer ${cookies.accessToken}`, }
    }
    ).then(res => {
      setMsg(res.data)
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

  return (
    <div>
      {!(cookies.accessToken === undefined || cookies.accessToken === "") && (
        <>
          <h1>ログインに成功しました。</h1>
          <button onClick={getCustomer}>getCustomer</button>
          <p>{msg}</p>
        </>
      )}
    </div>
  )
}

