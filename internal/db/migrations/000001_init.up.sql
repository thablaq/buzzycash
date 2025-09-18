CREATE TABLE public.admins (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(255),
    email character varying(255),
    password character varying(255),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    profile_picture character varying(255),
    role_id uuid
);
CREATE TABLE public.blacklisted_tokens (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    token text NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone
);
CREATE TABLE public.game_histories (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    ticket_type_id uuid,
    user_id uuid,
    transaction_history_id uuid,
    prize numeric,
    status text,
    winning_balls bigint,
    played_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.notifications (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    title character varying(255),
    message text,
    type text,
    is_read boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    subtitle character varying(500),
    amount integer,
    currency character varying(10),
    status character varying(50)
);
CREATE TABLE public.referral_earnings (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    wallet_id uuid NOT NULL,
    referrer_id uuid NOT NULL,
    referred_id uuid NOT NULL,
    points bigint NOT NULL,
    created_at timestamp with time zone,
    expires_at timestamp with time zone,
    used boolean DEFAULT false
);
CREATE TABLE public.referral_wallets (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    referral_balance bigint DEFAULT 0,
    points_used bigint DEFAULT 0,
    points_expired bigint DEFAULT 0,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    referrer_id uuid,
    referred_user_id uuid,
    points_earned numeric(10,2) DEFAULT 0,
    signup_date timestamp with time zone DEFAULT now(),
    first_transaction_date timestamp with time zone,
    transaction_count bigint DEFAULT 0,
    expires_at timestamp with time zone
);
CREATE TABLE public.referrals (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    referrer_id uuid,
    referred_user_id uuid,
    points_earned numeric(10,2) DEFAULT 0,
    points_used numeric(10,2) DEFAULT 0,
    points_expired numeric(10,2) DEFAULT 0,
    signup_date timestamp with time zone DEFAULT now(),
    first_transaction_date timestamp with time zone,
    transaction_count bigint DEFAULT 0,
    created_at timestamp with time zone,
    expires_at timestamp with time zone
);
CREATE TABLE public.refresh_tokens (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    token character varying(255),
    expire_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone
);
CREATE TABLE public.roles (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(255),
    description character varying(255)
);
CREATE TABLE public.ticket_purchases (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    total_amount numeric,
    unit_price numeric,
    quantity bigint,
    purchased_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    currency character varying(3) DEFAULT 'NGN'::character varying
);
CREATE TABLE public.transactions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    amount bigint,
    transaction_reference character varying(255),
    reference character varying(255),
    metadata jsonb,
    customer_email character varying(255),
    payment_status text,
    payment_type text,
    total_amount bigint,
    unit_price bigint,
    quantity bigint,
    currency text DEFAULT 'NGN'::text,
    paid_at timestamp with time zone,
    deleted_at timestamp with time zone,
    transaction_type text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone,
    payment_method text,
    category text
);
CREATE TABLE public.user_otp_securities (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    code character varying(255) NOT NULL,
    created_at timestamp with time zone,
    expires_at timestamp with time zone,
    retry_count bigint DEFAULT 0,
    locked_until timestamp with time zone,
    is_otp_verified_for_password_reset boolean DEFAULT false,
    sent_to character varying(255),
    action character varying(50) NOT NULL
);
CREATE TABLE public.user_otp_security (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    code character varying(255) NOT NULL,
    created_at timestamp with time zone,
    expires_at timestamp with time zone,
    retry_count bigint DEFAULT 0,
    locked_until timestamp with time zone,
    is_otp_verified_for_password_reset boolean DEFAULT false,
    sent_to character varying(255),
    action character varying(50) NOT NULL
);
CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    full_name character varying(255),
    phone_number character varying(255),
    email character varying(255),
    username character varying(255),
    date_of_birth character varying(255),
    password character varying(255),
    profile_picture character varying(255),
    is_profile_created boolean DEFAULT false,
    referral_code character varying(255),
    is_active boolean DEFAULT true,
    is_email_verified boolean DEFAULT false,
    is_verified boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    last_login timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    referred_by_id uuid,
    gender text,
    country_of_residence character varying(255)
);
CREATE TABLE public.withdrawal_requests (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    amount numeric,
    currency character varying(3) DEFAULT 'NGN'::character varying,
    payment_status text,
    reason character varying(255),
    payment_reference character varying(255),
    requested_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    processed_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone,
    transaction_history_id uuid
);
ALTER TABLE ONLY public.admins
ALTER TABLE ONLY public.blacklisted_tokens
ALTER TABLE ONLY public.game_histories
ALTER TABLE ONLY public.notifications
ALTER TABLE ONLY public.referral_earnings
ALTER TABLE ONLY public.referral_wallets
ALTER TABLE ONLY public.referrals
ALTER TABLE ONLY public.refresh_tokens
ALTER TABLE ONLY public.roles
ALTER TABLE ONLY public.ticket_purchases
ALTER TABLE ONLY public.transactions
ALTER TABLE ONLY public.blacklisted_tokens
ALTER TABLE ONLY public.user_otp_security
ALTER TABLE ONLY public.user_otp_securities
ALTER TABLE ONLY public.user_otp_security
ALTER TABLE ONLY public.users
ALTER TABLE ONLY public.withdrawal_requests
CREATE INDEX idx_referral_earnings_wallet_id ON public.referral_earnings USING btree (wallet_id);
ALTER TABLE ONLY public.referral_earnings
ALTER TABLE ONLY public.admins
ALTER TABLE ONLY public.ticket_purchases
ALTER TABLE ONLY public.game_histories
ALTER TABLE ONLY public.notifications
ALTER TABLE ONLY public.user_otp_security
ALTER TABLE ONLY public.user_otp_securities
ALTER TABLE ONLY public.referral_earnings
ALTER TABLE ONLY public.referral_wallets
ALTER TABLE ONLY public.referrals
ALTER TABLE ONLY public.referral_wallets
ALTER TABLE ONLY public.referrals
ALTER TABLE ONLY public.refresh_tokens
ALTER TABLE ONLY public.ticket_purchases
ALTER TABLE ONLY public.transactions
ALTER TABLE ONLY public.withdrawal_requests
