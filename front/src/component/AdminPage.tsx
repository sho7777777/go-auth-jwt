import { useEffect } from 'react';
import { Routes, Route, useNavigate } from 'react-router-dom';
import axios from "axios"

const AdminPage = (props: any) => {
  const { cookies } = props
  const navigate = useNavigate()

  useEffect(() => {
    axios.get("http://localhost:8000/admin", {
      headers: { Authorization: `Bearer ${cookies.accessToken}`, }
    }
      // 成功した場合（管理者の場合）はそのままページにとどまり、失敗した場合はトップに戻る
    ).then().catch(err => {
      navigate("/")
    })
  }, [])

  const onClickBack = () => {
    navigate("/resource")
  }

  return (
    <Routes>
      <Route path="/" element={
        <>
          <h1>管理者ページ</h1>
          <button onClick={onClickBack}>Resource Pageに戻る</button>
        </>
      }></Route>
    </Routes>
  )
}

export default AdminPage
