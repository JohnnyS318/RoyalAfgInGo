import React, { FC, useState } from "react";
import CurrencyInput from "react-currency-input-field";
import Layout from "../../components/layout";
import { useRouter } from "next/router";
import { getCSRF } from "../../hooks/auth/csrf";
import { useSnackbar } from "notistack";
import { GetServerSideProps } from "next";

export const getServerSideProps: GetServerSideProps = async (context) => {
    const csrf = await getCSRF(context);
    return {
        props: {
            csrf: csrf
        }
    };
};

type DepositProps = {
    csrf: string;
};

const Deposit: FC<DepositProps> = ({ csrf }) => {
    const { locale } = useRouter();
    const [loading, setLoading] = useState(false);
    const [success, setSuccess] = useState(false);
    const [amount, setAmount] = useState(-1);
    const { enqueueSnackbar } = useSnackbar();
    const router = useRouter();

    const onSubmit = async (e) => {
        e.preventDefault();
        if (amount <= 0) {
            return;
        }
        setLoading(true);
        const res = await fetch("/api/bank/deposit", {
            method: "POST",
            headers: {
                "X-CSRF-Token": csrf
            },
            body: JSON.stringify({
                amount: amount * 100
            })
        });

        if (!res.ok) {
            setLoading(false);
            setSuccess(false);
            console.log("Error during response: ", res);
            enqueueSnackbar("Error during depositing", { variant: "error" });
            return;
        }

        setLoading(false);
        setSuccess(true);
        enqueueSnackbar("Deposited successfully", { variant: "success" });
        router.push("/wallet").then();
    };

    return (
        <Layout disableFooter={true} headerAbsolute={true}>
            <div className="m-0 flex justify-center">
                {amount == 0 && (
                    <span className="text-sm bg-red-700 text-white px-3 py-1 rounded h-auto absolute top-16">
                        Specify a amount between $0.01 and 99,999,999.99
                    </span>
                )}
                <div className="font-sans text-5xl m-0 font-semibold grid justify-center items-center h-screen">
                    {!loading && !success && (
                        <form onSubmit={onSubmit}>
                            <div className="grid justify-center items-center">
                                <label htmlFor="amount" className="flex align-middle items-center justify-center">
                                    Deposit
                                    <CurrencyInput
                                        name="amount"
                                        className="flex w-80 text-center font-sans font-semibold border-black placeholder-blue-300 text-blue-600 border-solid border-b-2 mx-4 outline-none"
                                        placeholder="amount"
                                        intlConfig={{ locale: locale, currency: "USD" }}
                                        autoComplete="off"
                                        defaultValue={0}
                                        onValueChange={(val) => {
                                            setAmount(parseFloat(val));
                                        }}
                                        allowNegativeValue={false}
                                    />
                                    to your account.
                                </label>
                                <button
                                    className="text-lg bg-black hover:bg-blue-600 transition-colors duration-200 ease-in-out flex text-center px-6 py-1 w-fit text-white rounded-md mt-12 justify-self-center"
                                    type="submit">
                                    Continue
                                </button>
                            </div>
                        </form>
                    )}
                    {!loading && success && <h1>Success</h1>}
                    {loading && <h1>Loading...</h1>}
                </div>
            </div>
        </Layout>
    );
};

export default Deposit;
