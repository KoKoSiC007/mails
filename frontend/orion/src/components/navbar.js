import { isLogged } from "../internal/utils"

export const Navbar = () => {
    const logged = isLogged()

    return (
        <nav>
            <ul className="menu">
                <li key='home'>
                    <a href='/'>Home</a>
                </li>
                <li key='about'>
                    <a href='/about'>About</a>
                </li>
                <li key='mails'>
                    <a href='/mails'>Mails</a>
                </li>

                {logged ? null : <li key='auth'><a href='/login'>Login</a></li>}
                {logged ? <li key='currencies'><a href='/currencies'>Currencies</a></li> : null}
            </ul>
        </nav>
    )
}