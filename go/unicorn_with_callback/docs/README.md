# Brianterry::Unicorn::Maker

A resource that creates unicorns.

## Syntax

To declare this entity in your AWS CloudFormation template, use the following syntax:

### JSON

<pre>
{
    "Type" : "Brianterry::Unicorn::Maker",
    "Properties" : {
        "<a href="#name" title="Name">Name</a>" : <i>String</i>,
        "<a href="#color" title="Color">Color</a>" : <i>String</i>
    }
}
</pre>

### YAML

<pre>
Type: Brianterry::Unicorn::Maker
Properties:
    <a href="#name" title="Name">Name</a>: <i>String</i>
    <a href="#color" title="Color">Color</a>: <i>String</i>
</pre>

## Properties

#### Name

The name of the majestic animal

_Required_: Yes

_Type_: String

_Minimum_: <code>3</code>

_Maximum_: <code>250</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Color

The Color of the majestic animal

_Required_: Yes

_Type_: String

_Minimum_: <code>3</code>

_Maximum_: <code>250</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

## Return Values

### Ref

When you pass the logical ID of this resource to the intrinsic `Ref` function, Ref returns the UID.

### Fn::GetAtt

The `Fn::GetAtt` intrinsic function returns a value for a specified attribute of this type. The following are the available attributes and sample return values.

For more information about using the `Fn::GetAtt` intrinsic function, see [Fn::GetAtt](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-getatt.html).

#### UID

The ID of the majestic animal

