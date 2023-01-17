import { useEffect, useState } from 'react';
import { Routes, Route, useNavigate } from 'react-router-dom';
import axios from "axios"

const AdminPage = (props: any) => {
  const { cookies } = props
  const navigate = useNavigate()
  // レスポンスの速度によってコンテンツが表示されてしまうのを防ぐため、最初はコンテンツを表示しない
  const [content, setContent] = useState(false)

  // 直アクセス防止
  useEffect(() => {
    axios.get("http://localhost:8000/admin", {
      headers: { Authorization: `Bearer ${cookies.accessToken}`, }
    }
      // 成功した場合（管理者の場合）はそのままページにとどまり、失敗した場合はトップに戻る
    ).then(res => {
      setContent(true)
    }).catch(err => {
      navigate("/")
    })
  }, [])

  const onClickBack = () => {
    navigate("/resource")
  }

  return (
    <Routes>
      {content && (
        <Route path="/" element={
          <>
            <h1>管理者ページ</h1>
            <button onClick={onClickBack}>Resource Pageに戻る</button>
          </>
        }></Route>
      )}
    </Routes>
  )
}

export default AdminPage
