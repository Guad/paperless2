import React, { FunctionComponent } from 'react';
import PropTypes from 'prop-types';
import { Field, Form, FormSpy } from 'react-final-form';
import CardActions from '@material-ui/core/CardActions';
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';
import CircularProgress from '@material-ui/core/CircularProgress';
import { makeStyles, Theme } from '@material-ui/core/styles';
import { useTranslate, useLogin, useNotify, useSafeSetState } from 'ra-core';
import { useHistory } from 'react-router-dom';


const useStyles = makeStyles(
    (theme) => ({
        form: {
            padding: '0 1em 1em 1em',
        },
        input: {
            marginTop: '1em',
        },
        button: {
            width: '100%',
        },
        icon: {
            margin: '10px auto',
        },
    }),
    { name: 'RaLoginForm' }
);

const Input = ({
    meta: { touched, error }, // eslint-disable-line react/prop-types
    input: inputProps, // eslint-disable-line react/prop-types
    ...props
}) => (
        <TextField
            error={!!(touched && error)}
            helperText={touched && error}
            {...inputProps}
            {...props}
            fullWidth
        />
    );

const LoginForm = ({ redirectTo }) => {
    const [loading, setLoading] = useSafeSetState(false);
    const login = useLogin();
    const translate = useTranslate();
    const notify = useNotify();
    const classes = useStyles({});
    const history = useHistory();

    const validate = (values) => {
        const errors = { username: undefined, password: undefined };

        if (!values.username) {
            errors.username = translate('ra.validation.required');
        }
        if (!values.password) {
            errors.password = translate('ra.validation.required');
        }
        return errors;
    };

    const submit = values => {
        console.log(values);
        setLoading(true);
        login(values, redirectTo)
            .then(() => {
                setLoading(false);
            })
            .catch(error => {
                setLoading(false);
                notify(
                    typeof error === 'string'
                        ? error
                        : typeof error === 'undefined' || !error.message
                            ? 'ra.auth.sign_in_error'
                            : error.message,
                    'warning'
                );
            });
    };

    const register = values => {
        if (!values.username || !values.password) {
            notify('Username or password are empty.', 'warning');
            return;
        }

        setLoading(true);
        fetch('/api/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(values),
        }).then((resp) => {
            if (!resp.ok) throw new Error('Email is already registered!');

            return resp.json();
        }).then((data) => {
            setLoading(false);
            history.push('/');
        }).catch(error => {
            notify(
                typeof error === 'string'
                    ? error
                    : typeof error === 'undefined' || !error.message
                        ? 'ra.auth.sign_in_error'
                        : error.message,
                'warning'
            );
            setLoading(false);
        })
    };

    return (
        <Form
            onSubmit={submit}
            validate={validate}
            render={({ handleSubmit }) => (
                <form onSubmit={handleSubmit} noValidate>
                    <div className={classes.form}>
                        <div className={classes.input}>
                            <Field
                                autoFocus
                                id="username"
                                name="username"
                                component={Input}
                                label={translate('ra.auth.username')}
                                disabled={loading}
                            />
                        </div>
                        <div className={classes.input}>
                            <Field
                                id="password"
                                name="password"
                                component={Input}
                                label={translate('ra.auth.password')}
                                type="password"
                                disabled={loading}
                                autoComplete="current-password"
                            />
                        </div>
                        <Field component={'input'} id="register" name="register" type="hidden" />
                    </div>
                    {loading ? (
                        <CardActions>                            
                            <CircularProgress
                                className={classes.icon}
                                size={24}
                                thickness={3}
                            />
                        </CardActions>
                    ) : (
                            <CardActions>
                                <Button
                                    variant="contained"
                                    type="submit"
                                    color="primary"
                                    disabled={loading}
                                    className={classes.button}
                                >
                                    {translate('ra.auth.sign_in')}
                                </Button>
                                <FormSpy>
                                    {props => (
                                        <Button
                                            variant="contained"
                                            color="primary"
                                            disabled={loading}
                                            className={classes.button}
                                            onClick={_ => { register(props.form.getState().values) }}
                                        >
                                            Register
                                        </Button>
                                    )}
                                </FormSpy>
                            </CardActions>
                        )}
                </form>
            )}
        />
    );
};

LoginForm.propTypes = {
    redirectTo: PropTypes.string,
};

export default LoginForm;
