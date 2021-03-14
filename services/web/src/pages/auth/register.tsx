import React, { FC } from "react";
import { useForm } from "react-hook-form";
import FormItem from "../..//components/form/form-item";
import Layout from "../../components/layout";
import { register as registerAccount } from "../../hooks/auth";
import Head from "next/head";
import { formatTitle } from "../../utils/title";
import PasswordBox from "../../components/form/passwordBox";
import { useSnackbar } from "notistack";
import Checkbox from "@material-ui/core/Checkbox";
import { GetServerSideProps, InferGetServerSidePropsType } from "next";
import { getCSRF } from "../../hooks/auth/csrf";

type RegisterDto = {
    username: string;
    password: string;
    birthdate: Date;
    email: string;
    fullName: string;
};

export const getServerSideProps: GetServerSideProps = async (context) => {
    const csrf = await getCSRF(context);
    return {
        props: {
            csrf: csrf
        }
    };
};

const Register: FC = ({ csrf }: InferGetServerSidePropsType<typeof getServerSideProps>) => {
    const { register, handleSubmit, errors } = useForm<RegisterDto>();
    const { enqueueSnackbar } = useSnackbar();

    const onSubmit = async (data) => {
        console.log("Register");
        const res = await registerAccount(
            {
                username: data.username,
                password: data.password,
                email: data.email,
                birthdate: new Date(data.birthdate).getTime() / 1000,
                fullName: data.fullName
            },
            csrf
        );
        if (res.ok) {
            enqueueSnackbar("Successfully Registered", { variant: "success" });
            if (typeof window !== undefined) {
                window.location.href = "/";
            }
        }
    };

    return (
        <>
            <Head>
                <title>{formatTitle("Register a new Account")}</title>
            </Head>
            <Layout disableFooter>
                <div className="w-full md:h-screen flex items-center justify-center md:absolute md:inset-0">
                    <div className="bg-gray-200 md:rounded-md shadow-md">
                        <div className="heading mx-16 mt-8 mb-2">
                            <h1 className="text-center font-sans font-semibold text-3xl">Register A New Account</h1>
                        </div>
                        <div className="content md:px-24 px-4">
                            <form onSubmit={handleSubmit(onSubmit)}>
                                <FormItem>
                                    <label htmlFor="username" className="mb-2 block">
                                        Username*:
                                    </label>
                                    <input
                                        className="block px-4 py-4 rounded w-full"
                                        ref={register({ required: true, maxLength: 100, minLength: 3 })}
                                        type="text"
                                        id="username"
                                        name="username"
                                        placeholder="Your Username"
                                    />
                                    {errors.username && (
                                        <span className="text-sm text-red-700">
                                            This field is required and can only be more than 3 and less than 100!
                                        </span>
                                    )}
                                </FormItem>
                                <PasswordBox errors={errors} register={register} />
                                <FormItem>
                                    <label htmlFor="birthdate" className="mb-2 block">
                                        Birthdate*:
                                    </label>
                                    <input
                                        className="block px-4 py-4 rounded w-full"
                                        ref={register({ required: true })}
                                        type="date"
                                        id="birthdate"
                                        name="birthdate"
                                    />
                                    {errors.birthdate && <span className="text-sm text-red-700">This field is required!</span>}
                                </FormItem>
                                <FormItem>
                                    <label htmlFor="email">Email*:</label>
                                    <input
                                        className="block px-4 py-4 rounded w-full"
                                        ref={register({
                                            required: true,
                                            minLength: "3",
                                            maxLength: "100"
                                        })}
                                        name="email"
                                        id="email"
                                        type="email"
                                        placeholder="Your Email"
                                    />
                                    {errors.birthdate && <span className="text-sm text-red-700">This field is required</span>}
                                </FormItem>
                                <FormItem>
                                    <label htmlFor="fullname">Fullname*:</label>
                                    <input
                                        className="block px-4 py-4 rounded w-full"
                                        ref={register({
                                            required: true,
                                            minLength: "3",
                                            maxLength: "100"
                                        })}
                                        name="fullName"
                                        id="fullName"
                                        type="fullName"
                                        placeholder="Your Name"
                                    />
                                    {errors.birthdate && <span className="text-sm text-red-700">This field is required</span>}
                                </FormItem>
                                <div className="mb-4 font-sans text-lg font-medium">
                                    <Checkbox value="on" color="primary" required></Checkbox>
                                    <span>
                                        I consent to the{" "}
                                        <a href="/legal/terms" className="font-sans text-blue-800">
                                            terms and conditions
                                        </a>{" "}
                                        and our{" "}
                                        <a href="/legal/privacy" className="font-sans text-blue-800">
                                            privacy statement
                                        </a>
                                    </span>
                                </div>
                                <button
                                    className="block w-full px-4 py-2  bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors duration-150 font-sans font-medium cursor-pointer"
                                    type="submit">
                                    Register
                                </button>
                                <span className="font-sans font-light text-sm mb-8">
                                    Or{" "}
                                    <a href="/login" className="font-sans text-blue-800">
                                        login
                                    </a>{" "}
                                    if you already have an account
                                </span>
                                <span className="text-sm mb-8 font-sans font-light block text-ce">Textfields with a * are required</span>
                            </form>
                        </div>
                    </div>
                </div>
            </Layout>
        </>
    );
};

export default Register;
