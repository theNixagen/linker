"use client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { zodResolver } from "@hookform/resolvers/zod";
import { SubmitHandler, useForm } from "react-hook-form";
import z from "zod";
import { createUser } from "./actions";

const schema = z.object({
  name: z.string({ message: "Nome é obrigatório" }),
  email: z
    .string({ message: "Email é obrigatório" })
    .email({ message: "email invalido" }),
  password: z
    .string({ message: "Senha é obrigatória" })
    .min(6, { message: "A senha deve ter pelo menos 6 caracteres" }),
});

type FormValues = z.infer<typeof schema>;

export default function Register() {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FormValues>({ resolver: zodResolver(schema) });

  const onSubmit: SubmitHandler<FormValues> = async (data) => {
    await createUser(data);
  };

  return (
    <div className="text-black min-h-screen bg-gradient-to-br from-slate-900 to-slate-800 flex items-center justify-center px-4 py-10">
      <div className="w-full max-w-md rounded-3xl bg-white/95 p-8 shadow-2xl border border-white/30">
        <div className="mb-8 text-center space-y-2">
          <p className="text-sm uppercase tracking-[0.2em] text-slate-400">
            Linker
          </p>
          <div>
            <h1 className="text-3xl font-semibold text-slate-900">Registrar</h1>
            <p className="text-sm text-slate-500">
              Preencha seus dados para entrar
            </p>
          </div>
        </div>

        <form className="space-y-5" onSubmit={handleSubmit(onSubmit)}>
          <div className="space-y-2">
            <label
              className="text-sm font-medium text-slate-700"
              htmlFor="name"
            >
              Nome
            </label>
            <Input
              id="name"
              placeholder="Digite seu nome"
              type="text"
              {...register("name")}
            />
            {errors.name && (
              <p className="text-sm text-red-500">{errors.name.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <label
              className="text-sm font-medium text-slate-700"
              htmlFor="email"
            >
              Email
            </label>
            <Input
              id="email"
              placeholder="voce@email.com"
              type="email"
              {...register("email")}
            />
            {errors.email && (
              <p className="text-sm text-red-500">{errors.email.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <label
              className="text-sm font-medium text-slate-700"
              htmlFor="password"
            >
              Senha
            </label>
            <Input
              id="password"
              placeholder="Digite sua senha"
              type="password"
              {...register("password")}
            />
            {errors.password && (
              <p className="text-sm text-red-500">{errors.password.message}</p>
            )}
          </div>

          <Button
            type="submit"
            className="w-full bg-lime-400 text-black hover:bg-lime-500"
          >
            Registrar
          </Button>
        </form>
      </div>
    </div>
  );
}
