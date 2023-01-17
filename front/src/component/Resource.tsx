import axios from "axios"
import { useEffect, useState } from 'react';
import AdminPage from "./AdminPage";
import { Routes, Route, useNavigate } from 'react-router-dom';

export const Resource = (props: any) => {
  const { cookies, setCookies } = props
  const [resource, setResouce] = useState("")
  const [msg] = useState("")
  // レスポンスの速度によってコンテンツが表示されてしまうのを防ぐため、最初はコンテンツを表示しない
  const [content, setContent] = useState(false)
  const navigate = useNavigate()

  // 直アクセス防止
  useEffect(() => {
    console.log("useEffect")
    axios.get("http://localhost:8000/resource/1234", {
      headers: { Authorization: `Bearer ${cookies.accessToken}`, }
    }
      // 成功した場合（管理者の場合）はそのままページにとどまり、失敗した場合はトップに戻る
    ).then(res => {
      console.log("コンテンツ取得に成功")
      setContent(true)
    }).catch(err => {
      navigate("/")
    })
    // リフレッシュ用
  }, [cookies.accessToken])


  const getCustomer = () => {
    axios.get("http://localhost:8000/resource/1234", {
      headers: { Authorization: `Bearer ${cookies.accessToken}`, }
    }
      // 成功した場合は取得したリソースをstateにセットする
    ).then(res => {
      setResouce(res.data)
      // 失敗した場合、トークンをリフレッシュする
    }).catch(err => {
      axios.post("http://localhost:8001/auth/refresh", {
        access_token: cookies.accessToken,
        refresh_token: cookies.refreshToken,
        // リフレッシュに成功した場合は、新規に取得したアクセストークンをクッキーにセットする
      }).then(res => {
        setCookies("accessToken", res.data.access_token, { maxAge: 3600 })
        // リフレッシュに失敗した場合はトップページに遷移する
      }).catch(
        err => {
          navigate("/")
        })
    })
  }

  const goToAdminPage = () => {
    axios.get("http://localhost:8000/admin", {
      headers: { Authorization: `Bearer ${cookies.accessToken}`, }
    }
      // 成功した場合は管理者ページに遷移
    ).then(res => {
      navigate("admin")
      // 失敗した場合はトークンのリフレッシュを行う。
    }).catch(err => {
      axios.post("http://localhost:8001/auth/refresh", {
        access_token: cookies.accessToken,
        refresh_token: cookies.refreshToken,
        // 成功した場合は新規に取得したアクセストークンをクッキーにセットし、管理者ページに遷移する
      }).then(res => {
        setCookies("accessToken", res.data.access_token, { maxAge: 3600 })
        navigate("/resource/admin")
        // 失敗した場合はトップページに戻る
      }).catch(err => navigate("/"))
    })

  }

  return (
    <Routes>
      {content && (
        <>
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
        </>
      )}
    </Routes>
  )
}

