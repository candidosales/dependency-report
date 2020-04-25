export interface Repository {
    name: string;
    version: string;
    filter: string;
    url: string;
    documentation: Documentation;
    notifications: Array<GithubNotification>;
    updates: Array<Update>;
}

export interface Documentation {
    frontend: string;
    design: string;
}

export interface Update {
    name: string;
    current: string;
    update: string;
}


export interface GithubNotification {
    subject: GithubNotificationSubject;
    updated_at: Date;
    repository: GithubNotificationRepository;
}

export interface GithubNotificationSubject {
    title: string;
    url: string;
    latest_comment_url: string;
    type: string;
}

export interface GithubNotificationRepository {
    html_url: string;
}


