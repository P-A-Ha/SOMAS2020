import React from "react";
import logo from "../../../assets/logo/logo512.png";
import styles from "./IITO.module.css";

const IITO = () => {
  return (
    <div className={styles.root}>
      <img src={logo} className={styles.appLogo} alt="logo" />
      <p className={styles.text}>IITO Visualisation</p>
    </div>
  );
};

export default IITO;
