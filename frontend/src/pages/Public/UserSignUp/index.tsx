import { Formik, Form, FormikHelpers,Field ,FieldProps} from 'formik'
import * as yup from 'yup'
import { Button, Select, FormLabel, Heading } from '@chakra-ui/react'

import { signup } from 'api/auth'
import history from 'utils/history'

import { Input } from 'components/Form'

interface Values {
  username: string
  password: string
  user_type: string
}

const UserSignUp = () => {
  const initialValues: Values = { username: '', password: '', user_type: '' }

  const validationSchema: yup.SchemaOf<Values> = yup.object({
    username: yup.string().required('Required'),
    password: yup.string().required('Required'),
    user_type: yup.string().required('Required'),
  })

  const handleSignup = (values: Values, formikHelpers: FormikHelpers<Values>) => {
    console.log({ values, formikHelpers })
    signup({ username: values.username, password: values.password, user_type: values.user_type })
    formikHelpers.setSubmitting(false)
    history.push('/login')
  }

  return (
    <div>
      <Heading>USER SIGNUP</Heading>
      <Formik initialValues={initialValues} onSubmit={handleSignup} validationSchema={validationSchema}>
        <Form>
          <FormLabel htmlFor="username">Username</FormLabel>
          <Input id="username" name="username" placeholder="Username" />
          <FormLabel htmlFor="password">Password</FormLabel>
          <Input id="password" name="password" placeholder="Password" />

          <Field name="user_type"></Field>
          <Field>
            {({ field ,form}: FieldProps) => (
              <Select name="user_type"id="user_type" onChange={field.onChange}>
                <option>Select usertype</option>
                <option value="USER">User</option>
                <option value="ADMIN">Admin</option>
              </Select>
            )}
          </Field>
          <Button type="submit">Submit</Button>
        </Form>
      </Formik>
    </div>
  )
}

export default UserSignUp
