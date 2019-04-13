package resource

import (
	"errors"
	"fmt"
)

// getRoleID - gets a role_id from a role_name
func (r *Resource) getRoleID() error {
	if len(r.config.Source.RoleName) <= 0 {
		return errors.New("no role_name provided")
	}
	resp, err := r.client.Logical().Read(
		fmt.Sprintf(
			"auth/approle/role/%s/role-id",
			r.config.Source.RoleName,
		),
	)
	if err != nil {
		return err
	}

	if roleID, ok := resp.Data["role_id"]; ok {
		r.roleID = roleID.(string)
		r.logger.Debug().Msg("role_id success")
		return nil
	}

	return errors.New("no role_id returned")
}

// getSecretID - gets a secret_id for a role_id
func (r *Resource) getSecretID() error {
	if len(r.roleID) <= 0 {
		return errors.New("no role_id provided")
	}
	resp, err := r.client.Logical().Write(
		fmt.Sprintf(
			"auth/approle/role/%s/secret-id",
			r.config.Source.RoleName,
		),
		nil,
	)
	if err != nil {
		return err
	}

	if secretID, ok := resp.Data["secret_id"]; ok {
		r.secretID = secretID.(string)
		r.logger.Debug().Msg("secret_id success")
		return nil
	}

	return errors.New("no secret_id returned")
}

// loginWithAppRole - login via approle
func (r *Resource) loginWithAppRole() error {
	if len(r.roleID) <= 0 && len(r.secretID) <= 0 {
		return errors.New("role_id or secret_id not provided when authenticating")
	}

	resp, err := r.client.Logical().Write(
		"auth/approle/login",
		map[string]interface{}{
			"role_id":   r.roleID,
			"secret_id": r.secretID,
		},
	)
	if err != nil {
		return err
	}

	if resp.Auth == nil {
		return errors.New("no authentication returned")
	}

	r.client.SetToken(resp.Auth.ClientToken)
	return nil
}

// renewToken - renews a token
func (r *Resource) renewToken() error {
	if len(r.config.Source.VaultToken) <= 0 {
		return errors.New("error renewing vault client token, no vault_token provided")
	}

	r.logger.Debug().Msg("attempting renewal of token")

	resp, err := r.client.Logical().Write("auth/token/renew-self", nil)
	if err != nil {
		return err
	}
	if resp.Auth != nil {
		r.client.SetToken(resp.Auth.ClientToken)
		r.logger.Debug().Msg("succesfully renewed token")
		return nil
	}

	return errors.New("error renewing token")
}

// setToken - sets the vault client token
func (r *Resource) setToken() error {
	if len(r.config.Source.RoleName) > 0 {
		err := r.getRoleID()
		if err != nil {
			return err
		}

		err = r.getSecretID()
		if err != nil {
			return err
		}

		err = r.loginWithAppRole()
		if err != nil {
			return err
		}
		return nil
	}

	if len(r.config.Source.RoleID) > 0 && len(r.config.Source.SecretID) > 0 {
		r.roleID = r.config.Source.RoleID
		r.secretID = r.config.Source.SecretID
		err := r.loginWithAppRole()
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}
