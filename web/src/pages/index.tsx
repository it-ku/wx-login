import React, { useState, useEffect } from 'react';
import { Card, Input, Button, message, Spin, Typography } from 'antd';
import { login, getUser, logout, setToken, removeToken, UserInfo } from '../services/api';
import styles from './index.less';
import qrcode from '@/assets/qrcode.jpg';

const { Title, Text } = Typography;

export default function IndexPage() {
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [user, setUser] = useState<UserInfo | null>(null);
  const [code, setCode] = useState('');

  const checkLoginStat = async () => {
    try {
      setLoading(true);
      const res = await getUser();
      if (res.code === 0 && res.data) {
        setUser(res.data);
      } else {
        removeToken();
        setUser(null);
      }
    } catch (e) {
      removeToken();
      setUser(null);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    checkLoginStat();
  }, []);

  const handleLogin = async () => {
    if (!code || code.length !== 6) {
      message.warning('请输入6位验证码');
      return;
    }
    try {
      setSubmitting(true);
      const res = await login(code);
      if (res.code === 0 && res.data) {
        setToken(res.data.token);
        setUser(res.data.user);
        message.success('登录成功');
      } else {
        message.error(res.message || '验证码无效或已过期');
      }
    } catch (e) {
      message.error('登录请求失败，请检查网络');
    } finally {
      setSubmitting(false);
    }
  };

  const handleLogout = async () => {
    try {
      setLoading(true);
      await logout();
      removeToken();
      setUser(null);
      setCode('');
      message.success('已退出登录');
    } catch (e) {
      removeToken();
      setUser(null);
      setCode('');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className={styles.container}>
        <Spin size="large" />
      </div>
    );
  }

  return (
    <div className={styles.container}>
      {user ? (
        <Card className={styles.loginCard} bordered={false}>
          <Title level={3} className={styles.title}>欢迎回来</Title>
          <div className={styles.userInfo}>
            <p><Text type="secondary">OpenID:</Text> <Text strong>{user.open_id}</Text></p>
            <p><Text type="secondary">关注状态:</Text> <Text strong>{user.subscribed ? '已关注' : '未关注'}</Text></p>
            <p><Text type="secondary">注册时间:</Text> <Text>{new Date(user.created_at).toLocaleString()}</Text></p>
          </div>
          <Button type="primary" danger block onClick={handleLogout} className={styles.btn}>
            退出登录
          </Button>
        </Card>
      ) : (
        <Card className={styles.loginCard} bordered={false}>
          <Title level={3} className={styles.title}>微信公众平台登录</Title>
          <div className={styles.qrcodeWrapper}>
            <img src={qrcode} alt="公众号二维码" className={styles.qrcode} />
          </div>
          <p className={styles.hint}>请关注公众号，发送 <Text keyboard strong>验证码</Text> 获取登录验证码</p>
          <div className={styles.form}>
            <Input
              size="large"
              placeholder="请输入6位验证码"
              maxLength={6}
              value={code}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => setCode(e.target.value.replace(/\D/g, ''))}
              onPressEnter={handleLogin}
            />
            <Button
              type="primary"
              size="large"
              block
              loading={submitting}
              onClick={handleLogin}
              className={styles.btn}
            >
              登 录
            </Button>
          </div>
          <div className={styles.footer}>
            <Text type="secondary" className={styles.footerText}>验证码有效期为10分钟</Text>
          </div>
        </Card>
      )}
    </div>
  );
}
