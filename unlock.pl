#!/usr/bin/perl

use strict;
use warnings;

use Net::LDAP qw(LDAP_SUCCESS);
use Config::Simple;

my $configfile = $ENV{'HOME'} . "/.config/unlock.ini";
my $CONF = new Config::Simple($configfile) or die print Config::Simple->error();

################################################################################
# example config
#domaincontroller = 'dc.example.com'
#basedn           = 'DC=example,DC=com'
#binddn           = 'CN=admin,OU=superadmins,DC=example,DC=com'
#pw_cmd           = "secret-tool lookup user admin"
################################################################################

my $domaincontroller = $CONF->param('domaincontroller');
my $basedn           = $CONF->param('basedn');
my $binddn           = $CONF->param('binddn');
my $pw_cmd           = $CONF->param('pw_cmd');
my $bindpw           = `$pw_cmd`;

# connect
my $ad = Net::LDAP->new( 'ldaps://'.$domaincontroller, 	verify  => 'none') or die "$@";

# bind
my $res = $ad->bind($binddn, password => $bindpw);
unless ($res->code == LDAP_SUCCESS){
	print "Fail To bind to ldap: " . $res->message . "\n";
	exit 2;
}

my $rtv = 0;
foreach my $user (@ARGV) {

	# check if user exists
	$res = $ad->search(
			base   => $basedn,
			scope  => 'sub',
			filter => "(sAMAccountName=${user})",
			attrs  => ['*']
		);

	my @entries;
	if ($res->code != LDAP_SUCCESS) {
		print "failed to search user:" . $res->message . "\n";
		$rtv= $rtv + 3;
		next
	}	else {
		if ($res->count < 1) {
			print "user for pattern ${user} not found\n";
			$rtv= $rtv + 4;
			next;
		}

		@entries =$res->entries;
	}

	# loop over all entries
	foreach my $entry (@entries) {

		# check if user is locked
		if ($entry->get_value("lockoutTime") == 0) {
			print "account " . $entry->get_value("sAMAccountName") . " is not locked\n";
		} else {
			# unlock user
			$res = $ad->modify($entry->get_value("distinguishedName"), replace => { lockoutTime => '0'});
			if ($res->code != LDAP_SUCCESS) {
				print "unable to unlock" . $entry->get_value("sAMAccountName") . "\n";
			} else {
				print "unlock " . $entry->get_value("sAMAccountName") . "\n"
			}
		}
	}
}
exit $rtv;
